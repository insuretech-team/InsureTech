package handlers

// b2b_bulk_upload_handler.go
// ──────────────────────────
// POST /v1/b2b/employees/bulk-upload
//
// Accepts multipart/form-data with:
//   - file        : Excel (.xlsx) or CSV file
//   - business_id : organisation UUID (required; also read from X-Business-ID header)
//
// Parses rows, maps columns to CreateEmployeeRequest fields,
// calls CreateEmployee for each row via gRPC, and returns a
// structured JSON summary: { ok, message, result: { created, failed, total, errors[] } }
//
// Expected column headers (case-insensitive, any order):
//   name, employee_id, department_name, email, mobile_number,
//   date_of_birth, date_of_joining, gender, insurance_category,
//   coverage_amount, number_of_dependent
//
// department_name: matched case-insensitively against existing departments; created if not found.
// department_id is also accepted as a legacy alias and used directly if provided.
//
// The handler is intentionally lenient: unknown columns are ignored,
// missing optional fields use zero values.

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	b2bentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/grpc/metadata"
)

// ─── Bulk Upload Result ───────────────────────────────────────────────────────

type bulkUploadResult struct {
	Created int               `json:"created"`
	Failed  int               `json:"failed"`
	Total   int               `json:"total"`
	Errors  []bulkUploadError `json:"errors,omitempty"`
}

type bulkUploadError struct {
	Row     int    `json:"row"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message"`
}

// ─── Handler ─────────────────────────────────────────────────────────────────

// BulkUploadEmployees handles POST /v1/b2b/employees/bulk-upload
func (h *B2BServiceHandler) BulkUploadEmployees(w http.ResponseWriter, r *http.Request) {
	// 32 MB max upload
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSONError(w, http.StatusBadRequest, "could not parse multipart form: "+err.Error())
		return
	}

	// Business ID: form value → header (set by B2BContextMiddleware)
	businessID := strings.TrimSpace(r.FormValue("business_id"))
	if businessID == "" {
		businessID = strings.TrimSpace(r.Header.Get("X-Business-ID"))
	}
	if businessID == "" {
		writeJSONError(w, http.StatusBadRequest, "business_id is required")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "file field is required")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not read uploaded file")
		return
	}

	fileName := strings.ToLower(header.Filename)
	var rows [][]string

	switch {
	case strings.HasSuffix(fileName, ".xlsx"):
		rows, err = parseXLSX(fileBytes)
	case strings.HasSuffix(fileName, ".csv"):
		rows, err = parseCSV(fileBytes)
	default:
		// Try XLSX first, then CSV
		rows, err = parseXLSX(fileBytes)
		if err != nil {
			rows, err = parseCSV(fileBytes)
		}
	}
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "could not parse file: "+err.Error())
		return
	}

	if len(rows) < 2 {
		writeJSONError(w, http.StatusBadRequest, "file must have a header row and at least one data row")
		return
	}

	colMap := buildColMap(rows[0])
	ctx := withBulkMD(r.Context(), businessID, r.Header.Get("X-User-ID"), r.Header.Get("X-Portal"), r.Header.Get("X-Tenant-ID"))

	// Load plan catalog once — used to resolve plan_name → plan_id
	// planNameIndex: lowercase plan name → plan_id
	planNameIndex := make(map[string]string)
	if catalogResp, catErr := h.client.ListPurchaseOrderCatalog(ctx, &b2bservicev1.ListPurchaseOrderCatalogRequest{}); catErr == nil {
		for _, item := range catalogResp.GetItems() {
			if strings.TrimSpace(item.GetPlanName()) != "" {
				planNameIndex[strings.ToLower(strings.TrimSpace(item.GetPlanName()))] = item.GetPlanId()
			}
		}
	}

	// dept name → UUID cache to avoid repeated ListDepartments calls for same name
	deptNameCache := make(map[string]string)

	result := &bulkUploadResult{}
	dataRows := rows[1:]
	result.Total = len(dataRows)

	for i, row := range dataRows {
		rowNum := i + 2 // 1-based, header is row 1

		req, parseErr := rowToCreateEmployeeRequest(row, colMap, businessID)
		if parseErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, bulkUploadError{
				Row:     rowNum,
				Message: "parse error: " + parseErr.Error(),
			})
			continue
		}

		// Skip completely blank rows
		if strings.TrimSpace(req.Name) == "" && strings.TrimSpace(req.EmployeeId) == "" {
			result.Total--
			continue
		}

		if strings.TrimSpace(req.Name) == "" {
			result.Failed++
			result.Errors = append(result.Errors, bulkUploadError{
				Row:     rowNum,
				Message: "name is required",
			})
			continue
		}
		if strings.TrimSpace(req.EmployeeId) == "" {
			result.Failed++
			result.Errors = append(result.Errors, bulkUploadError{
				Row:     rowNum,
				Name:    req.Name,
				Message: "employee_id is required",
			})
			continue
		}

		// Resolve department_name → department_id if department_name column was used
		// and department_id was not already set directly.
		deptName := getCol(row, colMap, "department_name")
		if strings.TrimSpace(deptName) != "" && strings.TrimSpace(req.DepartmentId) == "" {
			// Validate the department name — reject values that are clearly not department names.
			// This prevents reference column bleed-over (e.g. premium values like "430", "500",
			// or column header text like "premium_amount (BDT)") from being created as departments.
			if err := validateDepartmentName(deptName); err != nil {
				result.Failed++
				result.Errors = append(result.Errors, bulkUploadError{
					Row:     rowNum,
					Name:    req.Name,
					Message: "invalid department_name: " + err.Error(),
				})
				continue
			}
			key := strings.ToLower(strings.TrimSpace(deptName))
			if cachedID, ok := deptNameCache[key]; ok {
				req.DepartmentId = cachedID
			} else {
				resolveCtx, resolveCancel := context.WithTimeout(ctx, 10*time.Second)
				deptID, resolveErr := h.ensureAdminDepartment(resolveCtx, businessID, deptName)
				resolveCancel()
				if resolveErr != nil {
					result.Failed++
					result.Errors = append(result.Errors, bulkUploadError{
						Row:     rowNum,
						Name:    req.Name,
						Message: "department resolve error: " + grpcErrMsg(resolveErr),
					})
					continue
				}
				deptNameCache[key] = deptID
				req.DepartmentId = deptID
			}
		}

		// Resolve assigned_plan_name → assigned_plan_id if plan name column was used
		// and plan_id was not already set directly.
		planName := getCol(row, colMap, "assigned_plan_name")
		if strings.TrimSpace(planName) != "" && strings.TrimSpace(req.AssignedPlanId) == "" {
			if planID, ok := planNameIndex[strings.ToLower(strings.TrimSpace(planName))]; ok {
				req.AssignedPlanId = planID
			}
			// If not found in catalog, leave blank — not a fatal error, plan assignment is optional
		}

		rowCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		_, grpcErr := h.client.CreateEmployee(rowCtx, req)
		cancel()

		if grpcErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, bulkUploadError{
				Row:     rowNum,
				Name:    req.Name,
				Message: grpcErrMsg(grpcErr),
			})
		} else {
			result.Created++
		}
	}

	status := http.StatusOK
	if result.Created == 0 && result.Failed > 0 {
		status = http.StatusUnprocessableEntity
	}

	var msg string
	switch {
	case result.Failed == 0 && result.Created > 0:
		msg = fmt.Sprintf("All %d employees uploaded successfully.", result.Created)
	case result.Created > 0 && result.Failed > 0:
		msg = fmt.Sprintf("%d of %d employees uploaded successfully. %d rows had errors and were skipped — see details below.", result.Created, result.Total, result.Failed)
	case result.Created == 0 && result.Failed > 0:
		msg = fmt.Sprintf("No employees were uploaded. All %d rows had errors — see details below.", result.Failed)
	default:
		msg = "No data rows found in file."
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":      result.Failed == 0,
		"message": msg,
		"result":  result,
	})
}

// ─── Column mapping ───────────────────────────────────────────────────────────

var colAliases = map[string]string{
	// Name
	"name": "name", "full name": "name", "employee name": "name", "fullname": "name",
	// Employee ID
	"employee_id": "employee_id", "employee id": "employee_id", "emp id": "employee_id",
	"emp_id": "employee_id", "employeeid": "employee_id", "staff id": "employee_id", "staff_id": "employee_id",
	// Department name (preferred — resolved to UUID by ensureDepartmentByName)
	"department_name": "department_name", "department name": "department_name",
	"dept name": "department_name", "dept_name": "department_name", "department": "department_name",
	// Department ID (legacy — used directly if provided)
	"department_id": "department_id", "department id": "department_id",
	"dept id": "department_id", "dept_id": "department_id", "departmentid": "department_id",
	// Contact
	"email": "email", "work email": "email",
	"mobile_number": "mobile_number", "mobile number": "mobile_number",
	"mobile": "mobile_number", "phone": "mobile_number", "phone_number": "mobile_number",
	// Dates
	"date_of_birth": "date_of_birth", "date of birth": "date_of_birth",
	"dob": "date_of_birth", "birth_date": "date_of_birth",
	"date_of_joining": "date_of_joining", "date of joining": "date_of_joining",
	"joining_date": "date_of_joining", "joining date": "date_of_joining",
	"doj": "date_of_joining", "start_date": "date_of_joining", "start date": "date_of_joining",
	// Gender
	"gender": "gender", "sex": "gender",
	// Insurance
	"insurance_category": "insurance_category", "insurance category": "insurance_category",
	"insurance_type": "insurance_category", "insurance type": "insurance_category",
	"coverage_amount": "coverage_amount", "coverage amount": "coverage_amount", "coverage": "coverage_amount",
	"number_of_dependent": "number_of_dependent", "number of dependents": "number_of_dependent",
	"dependents": "number_of_dependent", "no of dependents": "number_of_dependent",
	"number_of_dependents": "number_of_dependent",
	// Assigned plan — name (preferred, resolved to UUID via catalog) or UUID directly
	"assigned_plan_name": "assigned_plan_name", "assigned plan name": "assigned_plan_name",
	"plan_name": "assigned_plan_name", "plan name": "assigned_plan_name", "plan": "assigned_plan_name",
	// Legacy: direct UUID (still accepted)
	"assigned_plan_id": "assigned_plan_id", "assigned plan id": "assigned_plan_id",
	"plan_id": "assigned_plan_id", "plan id": "assigned_plan_id", "planid": "assigned_plan_id",
}

// normalizeHeader strips all non-alphanumeric characters (spaces, underscores,
// dashes, parentheses, brackets, unicode punctuation) and lowercases — making
// matching truly robust regardless of separator style, locale, or extra decorators.
// e.g. "Coverage Amount (BDT)" → "coverageamountbdt"
//
//	"assigned_plan_name"    → "assignedplanname"
//	"Date of Birth"         → "dateofbirth"
func normalizeHeader(h string) string {
	var sb strings.Builder
	for _, r := range strings.ToLower(h) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// colAliasesNorm is the alias map with all keys pre-normalized via normalizeHeader.
// Built once at init time from colAliases.
var colAliasesNorm map[string]string

func init() {
	colAliasesNorm = make(map[string]string, len(colAliases))
	for k, v := range colAliases {
		colAliasesNorm[normalizeHeader(k)] = v
	}
}

func buildColMap(headerRow []string) map[string]int {
	m := make(map[string]int, len(headerRow))
	for i, h := range headerRow {
		normalized := normalizeHeader(h)
		if normalized == "" {
			continue
		}
		if canonical, ok := colAliasesNorm[normalized]; ok {
			if _, exists := m[canonical]; !exists { // first occurrence wins
				m[canonical] = i
			}
		}
	}
	return m
}

func getCol(row []string, colMap map[string]int, field string) string {
	idx, ok := colMap[field]
	if !ok || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

// ─── Row → Proto ──────────────────────────────────────────────────────────────

func rowToCreateEmployeeRequest(row []string, colMap map[string]int, businessID string) (*b2bservicev1.CreateEmployeeRequest, error) {
	doj := normalizeDate(getCol(row, colMap, "date_of_joining"))
	if doj == "" {
		doj = time.Now().UTC().Format("2006-01-02")
	}

	req := &b2bservicev1.CreateEmployeeRequest{
		BusinessId:        businessID,
		Name:              getCol(row, colMap, "name"),
		EmployeeId:        getCol(row, colMap, "employee_id"),
		DepartmentId:      getCol(row, colMap, "department_id"),
		Email:             getCol(row, colMap, "email"),
		MobileNumber:      getCol(row, colMap, "mobile_number"),
		DateOfBirth:       normalizeDate(getCol(row, colMap, "date_of_birth")),
		DateOfJoining:     doj,
		Gender:            parseGender(getCol(row, colMap, "gender")),
		InsuranceCategory: parseInsuranceCategory(getCol(row, colMap, "insurance_category")),
		AssignedPlanId:    getCol(row, colMap, "assigned_plan_id"),
	}

	if covStr := getCol(row, colMap, "coverage_amount"); covStr != "" {
		cov, cvErr := strconv.ParseFloat(strings.ReplaceAll(covStr, ",", ""), 64)
		if cvErr == nil && cov > 0 {
			req.CoverageAmount = &commonv1.Money{
				Amount:        int64(math.Round(cov * 100)),
				Currency:      "BDT",
				DecimalAmount: cov,
			}
		}
	}

	if depStr := getCol(row, colMap, "number_of_dependent"); depStr != "" {
		if dep, depErr := strconv.Atoi(depStr); depErr == nil {
			req.NumberOfDependent = int32(dep)
		}
	}

	return req, nil
}

// ─── Type converters ──────────────────────────────────────────────────────────

func parseGender(s string) b2bentityv1.EmployeeGender {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "M", "MALE", "EMPLOYEE_GENDER_MALE":
		return b2bentityv1.EmployeeGender_EMPLOYEE_GENDER_MALE
	case "F", "FEMALE", "EMPLOYEE_GENDER_FEMALE":
		return b2bentityv1.EmployeeGender_EMPLOYEE_GENDER_FEMALE
	case "O", "OTHER", "EMPLOYEE_GENDER_OTHER":
		return b2bentityv1.EmployeeGender_EMPLOYEE_GENDER_OTHER
	}
	return b2bentityv1.EmployeeGender_EMPLOYEE_GENDER_UNSPECIFIED
}

func parseInsuranceCategory(s string) commonv1.InsuranceType {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "LIFE", "LIFE INSURANCE", "INSURANCE_TYPE_LIFE", "1":
		return commonv1.InsuranceType_INSURANCE_TYPE_LIFE
	case "HEALTH", "HEALTH INSURANCE", "INSURANCE_TYPE_HEALTH", "2":
		return commonv1.InsuranceType_INSURANCE_TYPE_HEALTH
	case "AUTO", "MOTOR", "INSURANCE_TYPE_AUTO", "3":
		return commonv1.InsuranceType_INSURANCE_TYPE_AUTO
	case "TRAVEL", "INSURANCE_TYPE_TRAVEL", "4":
		return commonv1.InsuranceType_INSURANCE_TYPE_TRAVEL
	case "FIRE", "INSURANCE_TYPE_FIRE", "5":
		return commonv1.InsuranceType_INSURANCE_TYPE_FIRE
	case "MARINE", "INSURANCE_TYPE_MARINE", "6":
		return commonv1.InsuranceType_INSURANCE_TYPE_MARINE
	case "PROPERTY", "INSURANCE_TYPE_PROPERTY", "7":
		return commonv1.InsuranceType_INSURANCE_TYPE_PROPERTY
	case "LIABILITY", "INSURANCE_TYPE_LIABILITY", "8":
		return commonv1.InsuranceType_INSURANCE_TYPE_LIABILITY
	}
	return commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED
}

// normalizeDate accepts common date formats and returns YYYY-MM-DD.
func normalizeDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	formats := []string{
		"2006-01-02", "02/01/2006", "01/02/2006", "2006/01/02",
		"02-01-2006", "01-02-2006", "2 Jan 2006", "2 January 2006",
		"Jan 2, 2006", "January 2, 2006", "02 Jan 2006",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return s // return as-is; the B2B service will validate
}

// validateDepartmentName rejects values that are clearly not valid department names.
// Prevents reference-column bleed-over from creating junk departments in the DB.
// Rules:
//   - Must not be a pure number (e.g. "430", "500.00")
//   - Must not contain "reference" (e.g. "assigned_plan_name (reference)")
//   - Must not be a known column header pattern (contains "amount", "premium", "coverage" with special chars)
//   - Must be at least 2 characters after trimming
//   - Must contain at least one letter
func validateDepartmentName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return fmt.Errorf("department name must be at least 2 characters, got %q", name)
	}

	// Must contain at least one letter (rejects pure numbers, decimals)
	hasLetter := false
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r > 127 {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return fmt.Errorf("department name must contain letters, got %q (looks like a number or symbol)", name)
	}

	// Reject names that contain "reference" — these are template reference column headers
	lower := strings.ToLower(name)
	if strings.Contains(lower, "reference") {
		return fmt.Errorf("department name %q looks like a template column header, not a department name", name)
	}

	// Reject known column header patterns (e.g. "premium_amount (BDT)", "coverage_amount")
	columnPatterns := []string{
		"premium_amount", "coverage_amount", "assigned_plan", "insurance_category",
		"employee_id", "date_of_birth", "date_of_joining", "mobile_number",
	}
	for _, pattern := range columnPatterns {
		if strings.Contains(lower, pattern) {
			return fmt.Errorf("department name %q looks like a column header, not a department name", name)
		}
	}

	return nil
}

func grpcErrMsg(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if idx := strings.Index(msg, "desc = "); idx >= 0 {
		return msg[idx+7:]
	}
	return msg
}

func withBulkMD(ctx context.Context, businessID, userID, portal, tenantID string) context.Context {
	pairs := []string{"x-business-id", businessID}
	if userID != "" {
		// B2B service authz interceptor requires x-user-id; also set x-caller-id for audit trails
		pairs = append(pairs, "x-user-id", userID, "x-caller-id", userID)
	}
	// x-portal is required by resolveAuthzDomain to pick the correct Casbin domain
	// (e.g. "PORTAL_B2B" → "b2b:{orgId}", "PORTAL_SYSTEM" → "system:root")
	if portal != "" {
		pairs = append(pairs, "x-portal", portal)
	}
	if tenantID != "" {
		pairs = append(pairs, "x-tenant-id", tenantID)
	}
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		md2 := md.Copy()
		for i := 0; i+1 < len(pairs); i += 2 {
			md2.Set(pairs[i], pairs[i+1])
		}
		return metadata.NewOutgoingContext(ctx, md2)
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...))
}

// ─── CSV Parser ───────────────────────────────────────────────────────────────

func parseCSV(data []byte) ([][]string, error) {
	// Strip UTF-8 BOM (\xEF\xBB\xBF) if present — added by template route so
	// Excel renders Bengali characters correctly, but csv.Reader must not see it.
	if bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		data = data[3:]
	}
	r := csv.NewReader(bytes.NewReader(data))
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1 // Allow variable number of fields per row (some editors strip trailing empty commas)
	return r.ReadAll()
}

// ─── Robust XLSX Parser ───────────────────────────────────────────────────────
// XLSX files are ZIP archives containing XML. We extract the shared strings
// table and the first worksheet, then build [][]string rows.
//
// Robustness notes:
//   - Shared strings: supports both plain <t> and rich-text <r><t> runs (Bengali/emoji/mixed)
//   - Date cells: Excel stores dates as numeric serials (type "n" or no type + style "d").
//     We detect numeric values in date-formatted cells and convert to YYYY-MM-DD.
//   - Content-agnostic: works identically for Bengali, English, Arabic, or any UTF-8 text.
//   - Unknown column headers are ignored; column order is irrelevant.

// xlsxSharedStrings parses the shared strings table.
// Each <si> entry may be a plain <t> or rich-text with multiple <r><t> runs.
type xlsxSharedStrings struct {
	Items []xlsxSI `xml:"si"`
}

type xlsxSI struct {
	T    string    `xml:"t"` // plain text
	Runs []xlsxRun `xml:"r"` // rich-text runs
}

// text returns the full string value, concatenating rich-text runs when present.
func (si xlsxSI) text() string {
	if len(si.Runs) > 0 {
		var sb strings.Builder
		for _, r := range si.Runs {
			sb.WriteString(r.T)
		}
		return sb.String()
	}
	return si.T
}

type xlsxRun struct {
	T string `xml:"t"`
}

type xlsxSheet struct {
	Rows []xlsxRow `xml:"sheetData>row"`
}

type xlsxRow struct {
	Cells []xlsxCell `xml:"c"`
}

type xlsxCell struct {
	R  string `xml:"r,attr"` // cell reference e.g. "A1"
	T  string `xml:"t,attr"` // type: "s"=shared string, "inlineStr", "b"=bool, "n"=number, "e"=error, ""=number/date
	S  string `xml:"s,attr"` // style index (used to detect date formatting)
	V  string `xml:"v"`      // value
	IS struct {
		T    string    `xml:"t"`
		Runs []xlsxRun `xml:"r"`
	} `xml:"is"` // inline string
}

// xlsxStyles holds just the numFmtId for each xf entry so we can detect date styles.
type xlsxStyles struct {
	CellXfs struct {
		Xf []struct {
			NumFmtId string `xml:"numFmtId,attr"`
		} `xml:"xf"`
	} `xml:"cellXfs"`
}

// isDateStyleIndex returns true if the given 0-based xf index is a date/time format.
// Excel built-in date numFmtIds: 14-22, 27-36, 45-47, 50-58 (locale-dependent).
func isDateNumFmtId(id int) bool {
	return (id >= 14 && id <= 22) ||
		(id >= 27 && id <= 36) ||
		(id >= 45 && id <= 47) ||
		(id >= 50 && id <= 58)
}

// excelSerialToDate converts an Excel date serial number to YYYY-MM-DD.
// Excel epoch: December 30, 1899 (with the Lotus 1-2-3 leap year bug for 1900).
func excelSerialToDate(serial float64) string {
	// Excel incorrectly treats 1900 as a leap year; adjust for serials > 59
	if serial >= 60 {
		serial--
	}
	epoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	days := int(serial)
	t := epoch.AddDate(0, 0, days)
	return t.Format("2006-01-02")
}

func parseXLSX(data []byte) ([][]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("not a valid xlsx file: %w", err)
	}

	// ── Read shared strings table ─────────────────────────────────────────────
	var sharedStrings []string
	for _, f := range zr.File {
		if f.Name == "xl/sharedStrings.xml" {
			rc, openErr := f.Open()
			if openErr != nil {
				return nil, openErr
			}
			var ss xlsxSharedStrings
			decErr := xml.NewDecoder(rc).Decode(&ss)
			rc.Close()
			if decErr != nil {
				return nil, fmt.Errorf("shared strings parse error: %w", decErr)
			}
			for _, item := range ss.Items {
				sharedStrings = append(sharedStrings, item.text())
			}
			break
		}
	}

	// ── Read styles to detect date-formatted cells ────────────────────────────
	// dateStyleSet maps xf index → true if it uses a date numFmt
	dateStyleSet := make(map[int]bool)
	for _, f := range zr.File {
		if f.Name == "xl/styles.xml" {
			rc, openErr := f.Open()
			if openErr != nil {
				break
			}
			var styles xlsxStyles
			decErr := xml.NewDecoder(rc).Decode(&styles)
			rc.Close()
			if decErr != nil {
				break
			}
			for i, xf := range styles.CellXfs.Xf {
				if numFmtId, convErr := strconv.Atoi(xf.NumFmtId); convErr == nil {
					if isDateNumFmtId(numFmtId) {
						dateStyleSet[i] = true
					}
				}
			}
			break
		}
	}

	// ── Read first worksheet ──────────────────────────────────────────────────
	var sheet xlsxSheet
	for _, f := range zr.File {
		if f.Name == "xl/worksheets/sheet1.xml" {
			rc, openErr := f.Open()
			if openErr != nil {
				return nil, openErr
			}
			decErr := xml.NewDecoder(rc).Decode(&sheet)
			rc.Close()
			if decErr != nil {
				return nil, fmt.Errorf("sheet parse error: %w", decErr)
			}
			break
		}
	}

	if len(sheet.Rows) == 0 {
		return nil, fmt.Errorf("worksheet is empty")
	}

	// ── Determine max column index ────────────────────────────────────────────
	maxCol := 0
	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
			col := xlsxColIndex(cell.R)
			if col+1 > maxCol {
				maxCol = col + 1
			}
		}
	}

	// ── Build [][]string ──────────────────────────────────────────────────────
	var result [][]string
	for _, row := range sheet.Rows {
		rowData := make([]string, maxCol)
		for _, cell := range row.Cells {
			col := xlsxColIndex(cell.R)
			if col < 0 || col >= maxCol {
				continue
			}

			var val string
			switch cell.T {
			case "s": // shared string index
				if idx, idxErr := strconv.Atoi(strings.TrimSpace(cell.V)); idxErr == nil && idx >= 0 && idx < len(sharedStrings) {
					val = sharedStrings[idx]
				}

			case "inlineStr":
				// Inline string — may also have rich-text runs
				if len(cell.IS.Runs) > 0 {
					var sb strings.Builder
					for _, r := range cell.IS.Runs {
						sb.WriteString(r.T)
					}
					val = sb.String()
				} else {
					val = cell.IS.T
				}

			case "b": // boolean
				if cell.V == "1" {
					val = "TRUE"
				} else {
					val = "FALSE"
				}

			case "e": // error
				val = ""

			default:
				// Numeric or date — check style to distinguish date from plain number
				if strings.TrimSpace(cell.V) == "" {
					val = ""
					break
				}
				styleIdx, styleErr := strconv.Atoi(strings.TrimSpace(cell.S))
				if styleErr == nil && dateStyleSet[styleIdx] {
					// Date serial → YYYY-MM-DD
					if serial, floatErr := strconv.ParseFloat(strings.TrimSpace(cell.V), 64); floatErr == nil {
						val = excelSerialToDate(serial)
					} else {
						val = strings.TrimSpace(cell.V)
					}
				} else {
					// Plain number or text stored as number
					val = strings.TrimSpace(cell.V)
				}
			}

			rowData[col] = strings.TrimSpace(val)
		}

		// Skip all-empty rows
		empty := true
		for _, v := range rowData {
			if v != "" {
				empty = false
				break
			}
		}
		if !empty {
			result = append(result, rowData)
		}
	}
	return result, nil
}

// xlsxColIndex converts an Excel column reference (e.g. "A1", "BC3") to a 0-based column index.
// Only the letter prefix is used; the row number is ignored.
func xlsxColIndex(cellRef string) int {
	col := 0
	for _, ch := range cellRef {
		if ch >= 'A' && ch <= 'Z' {
			col = col*26 + int(ch-'A'+1)
		} else {
			break
		}
	}
	return col - 1
}
