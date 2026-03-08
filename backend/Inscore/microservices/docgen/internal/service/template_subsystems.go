package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	excelTemplateBase64Key = "excel_template_base64"
	excelSheetKey          = "excel_sheet"
	excelHeaderRowKey      = "excel_header_row"
	excelMetaSheetKey      = "excel_meta_sheet"
	taxRateKey             = "tax_rate"
)

var excelHeaderAliases = map[string][]string{
	"description": {"description", "item", "item_name", "service", "particular", "details"},
	"quantity":    {"quantity", "qty", "qnt", "count"},
	"unit_price":  {"unit_price", "unitprice", "unit_cost", "rate", "price"},
	"amount":      {"amount", "line_total", "linetotal", "total"},
}

func enrichTemplatePayload(templateName string, payload map[string]any) error {
	alignTemplateDomainPayload(templateName, payload)
	alignBillingInvoicePayload(templateName, payload)
	if err := mergeExcelTemplateData(templateName, payload); err != nil {
		return err
	}
	normalizeTemplateTotals(templateName, payload)
	return nil
}

func alignTemplateDomainPayload(templateName string, payload map[string]any) {
	name := strings.ToLower(strings.TrimSpace(templateName))
	switch name {
	case "b2b_po":
		alignB2BPurchaseOrderPayload(payload)
	case "b2b_pi":
		alignB2BInvoicePayload(payload)
	case "b2c_pi":
		alignB2CInvoicePayload(payload)
	case "b2c_po":
		alignB2CServiceOrderPayload(payload)
	}
}

func alignBillingInvoicePayload(templateName string, payload map[string]any) {
	if templateKind(templateName) != "invoice" {
		return
	}

	flattenEnvelope(payload, "invoice")

	// Align with billing Invoice proto fields.
	setIfMissingString(payload, "reference_id", asString(payload["invoice_id"]))
	setIfMissingString(payload, "org_reference_id", asString(payload["business_id"]))
	setIfMissingString(payload, "invoice_number", firstNonEmpty(payload["invoice_number"], payload["invoice_id"]))
	setIfMissingString(payload, "issue_date", normalizeDateLike(asString(payload["issued_at"])))
	setIfMissingString(payload, "paid_date", normalizeDateLike(asString(payload["paid_at"])))
	setIfMissingString(payload, "payment_reference", asString(payload["payment_id"]))
	setIfMissingString(payload, "business_reference", asString(payload["business_id"]))
	setIfMissingString(payload, "invoice_status", asString(payload["status"]))
	setIfMissingString(payload, "policy_ids_display", joinStringSlice(toStringSlice(payload["policy_ids"]), ", "))

	if due := normalizeDateLike(asString(payload["due_date"])); due != "" {
		payload["due_date"] = due
	}
	if issued := normalizeDateLike(asString(payload["issue_date"])); issued != "" {
		payload["issue_date"] = issued
	}
	if paid := normalizeDateLike(asString(payload["paid_date"])); paid != "" {
		payload["paid_date"] = paid
	}

	if asString(payload["invoice_status"]) == "" && payload["status"] != nil {
		payload["invoice_status"] = asString(payload["status"])
	}

	amount, hasAmount := billingAmountDecimal(payload["amount"])
	if hasAmount {
		payload["amount_display"] = formatMoney(amount)
		if amountFromPayload(payload, "subtotal") == 0 {
			payload["subtotal"] = formatMoney(amount)
		}
		if amountFromPayload(payload, "total") == 0 {
			payload["total"] = formatMoney(amount)
		}
		if _, exists := payload["items"]; !exists {
			payload["items"] = []any{
				map[string]any{
					"description": "Insurance premium charge",
					"quantity":    "1",
					"unit_price":  formatMoney(amount),
					"amount":      formatMoney(amount),
				},
			}
		}
	}
}

func alignB2BPurchaseOrderPayload(payload map[string]any) {
	orgRaw, _ := payload["organisation"].(map[string]any)
	if orgRaw == nil {
		orgRaw, _ = payload["organization"].(map[string]any)
	}
	departmentRaw, _ := payload["department"].(map[string]any)
	productRaw, _ := payload["product"].(map[string]any)
	planRaw, _ := payload["plan"].(map[string]any)

	flattenEnvelope(payload, "purchase_order")
	flattenEnvelope(payload, "purchaseOrder")
	flattenEnvelope(payload, "organisation")
	flattenEnvelope(payload, "organization")
	flattenEnvelope(payload, "department")
	flattenEnvelope(payload, "product")
	flattenEnvelope(payload, "plan")

	setIfMissing(payload, "reference_id", asString(payload["purchase_order_id"]))
	setIfMissing(payload, "issue_date", normalizeDateLike(asString(payload["created_at"])))
	setIfMissing(payload, "status_label", humanizeEnum(asString(payload["status"]), "PURCHASE_ORDER_STATUS_"))

	setIfMissing(payload, "org_name", firstNonEmpty(
		payload["org_name"],
		payload["business_name"],
		payload["organisation_name"],
		payload["organization_name"],
		orgRaw["name"],
		payload["name"],
		payload["company_name"],
	))
	setIfMissing(payload, "org_code", firstNonEmpty(
		payload["org_code"],
		payload["business_code"],
		orgRaw["code"],
		payload["code"],
	))
	setIfMissing(payload, "org_email", firstNonEmpty(
		payload["org_email"],
		orgRaw["contact_email"],
		payload["contact_email"],
		payload["business_email"],
	))
	setIfMissing(payload, "org_contact", firstNonEmpty(
		payload["org_contact"],
		orgRaw["contact_phone"],
		payload["contact_phone"],
		payload["business_phone"],
		payload["phone"],
	))
	setIfMissing(payload, "org_address", firstNonEmpty(
		payload["org_address"],
		orgRaw["address"],
		payload["business_address"],
		payload["address"],
	))

	setIfMissing(payload, "business_reference", asString(payload["business_id"]))
	setIfMissing(payload, "department_reference", asString(payload["department_id"]))
	setIfMissing(payload, "product_reference", asString(payload["product_id"]))
	setIfMissing(payload, "plan_reference", asString(payload["plan_id"]))
	setIfMissing(payload, "department_name", firstNonEmpty(payload["department_name"], payload["dept_name"], departmentRaw["name"]))
	setIfMissing(payload, "product_name", firstNonEmpty(payload["product_name"], payload["product_title"], productRaw["name"]))
	setIfMissing(payload, "plan_name", firstNonEmpty(payload["plan_name"], payload["plan_title"], planRaw["name"]))

	setIfMissing(payload, "employee_count_display", asString(payload["employee_count"]))
	setIfMissing(payload, "dependent_count_display", asString(payload["number_of_dependents"]))
	setIfMissing(payload, "insurance_category_label", humanizeEnum(asString(payload["insurance_category"]), "INSURANCE_TYPE_"))

	coverage := moneyToDisplay(payload["coverage_amount"])
	if coverage != "" {
		setIfMissing(payload, "coverage_amount_display", coverage)
	}
	estimated := moneyToDisplay(payload["estimated_premium"])
	if estimated != "" {
		setIfMissing(payload, "estimated_premium_display", estimated)
		setIfMissing(payload, "subtotal", estimated)
		setIfMissing(payload, "total", estimated)
	}
	setIfMissing(payload, "tax", "0.00")

	if _, ok := payload["items"]; !ok {
		payload["items"] = []any{
			map[string]any{
				"description": valueOr(payload, "plan_name", "Insurance plan enrollment"),
				"quantity":    valueOr(payload, "employee_count", "1"),
				"unit_price":  valueOr(payload, "estimated_premium_display", "0.00"),
				"amount":      valueOr(payload, "estimated_premium_display", "0.00"),
			},
		}
	}
}

func alignB2BInvoicePayload(payload map[string]any) {
	invoiceRaw, _ := payload["invoice"].(map[string]any)
	orgRaw, _ := payload["organisation"].(map[string]any)
	if orgRaw == nil {
		orgRaw, _ = payload["organization"].(map[string]any)
	}

	flattenEnvelope(payload, "invoice")
	flattenEnvelope(payload, "organisation")
	flattenEnvelope(payload, "organization")

	setIfMissing(payload, "invoice_number", firstNonEmpty(payload["invoice_number"], invoiceRaw["invoice_number"], payload["invoice_id"]))
	setIfMissing(payload, "issue_date", normalizeDateLike(asString(payload["issued_at"])))
	setIfMissing(payload, "paid_date", normalizeDateLike(asString(payload["paid_at"])))
	setIfMissing(payload, "due_date", normalizeDateLike(asString(firstNonEmpty(payload["due_date"], invoiceRaw["due_date"]))))
	setIfMissing(payload, "business_reference", asString(payload["business_id"]))
	setIfMissing(payload, "payment_reference", asString(payload["payment_id"]))
	setIfMissing(payload, "policy_ids_display", joinStringSlice(toStringSlice(payload["policy_ids"]), ", "))
	setIfMissing(payload, "invoice_status", humanizeEnum(asString(payload["status"]), "INVOICE_STATUS_", "PAYMENT_STATUS_"))
	setIfMissing(payload, "org_name", firstNonEmpty(
		payload["org_name"],
		payload["business_name"],
		payload["organisation_name"],
		payload["organization_name"],
		orgRaw["name"],
		payload["name"],
		payload["company_name"],
	))
	setIfMissing(payload, "org_email", firstNonEmpty(payload["org_email"], orgRaw["contact_email"], payload["contact_email"], payload["business_email"]))
	setIfMissing(payload, "org_contact", firstNonEmpty(payload["org_contact"], orgRaw["contact_phone"], payload["contact_phone"], payload["business_phone"], payload["phone"]))
	setIfMissing(payload, "org_address", firstNonEmpty(payload["org_address"], orgRaw["address"], payload["business_address"], payload["address"]))
}

func alignB2CInvoicePayload(payload map[string]any) {
	paymentRaw, _ := payload["payment"].(map[string]any)
	policyRaw, _ := payload["policy"].(map[string]any)

	flattenEnvelope(payload, "payment")
	flattenEnvelope(payload, "policy")

	setIfMissing(payload, "invoice_number", asString(payload["transaction_id"]))
	setIfMissing(payload, "reference_id", asString(payload["payment_id"]))
	setIfMissing(payload, "policy_reference", asString(payload["policy_id"]))
	setIfMissing(payload, "policy_number", asString(payload["policy_number"]))
	setIfMissing(payload, "payment_status", humanizeEnum(firstNonEmpty(payload["payment_status"], paymentRaw["status"], payload["status"]), "PAYMENT_STATUS_", "INVOICE_STATUS_"))
	setIfMissing(payload, "payment_method", humanizeEnum(firstNonEmpty(payload["payment_method"], payload["method"], paymentRaw["method"]), "PAYMENT_METHOD_"))
	setIfMissing(payload, "issue_date", normalizeDateLike(asString(payload["initiated_at"])))
	setIfMissing(payload, "paid_date", normalizeDateLike(asString(payload["completed_at"])))
	setIfMissing(payload, "coverage_start_date", normalizeDateLike(asString(payload["start_date"])))
	setIfMissing(payload, "coverage_end_date", normalizeDateLike(asString(payload["end_date"])))
	setIfMissing(payload, "policy_status", humanizeEnum(firstNonEmpty(payload["policy_status"], policyRaw["status"]), "POLICY_STATUS_"))
	setIfMissing(payload, "customer_reference", firstNonEmpty(payload["customer_id"], payload["payer_id"]))
	setIfMissing(payload, "customer_name", firstNonEmpty(payload["customer_name"], payload["policy_holder_name"], payload["payer_name"]))
	setIfMissing(payload, "customer_email", firstNonEmpty(payload["customer_email"], payload["payer_email"]))
	setIfMissing(payload, "customer_phone", firstNonEmpty(payload["customer_phone"], payload["payer_phone"]))
	setIfMissing(payload, "sum_insured_display", moneyToDisplay(payload["sum_insured"]))

	if total := moneyToDisplay(payload["amount"]); total != "" {
		setIfMissing(payload, "subtotal", total)
		setIfMissing(payload, "total", total)
		setIfMissing(payload, "amount_display", total)
		if _, ok := payload["items"]; !ok {
			payload["items"] = []any{
				map[string]any{
					"description": "Policy premium payment",
					"quantity":    "1",
					"unit_price":  total,
					"amount":      total,
				},
			}
		}
	}
	setIfMissing(payload, "tax", "0.00")
}

func alignB2CServiceOrderPayload(payload map[string]any) {
	serviceRaw, _ := payload["policy_service_request"].(map[string]any)
	if serviceRaw == nil {
		serviceRaw, _ = payload["policyServiceRequest"].(map[string]any)
	}

	flattenEnvelope(payload, "policy_service_request")
	flattenEnvelope(payload, "policyServiceRequest")
	flattenEnvelope(payload, "request_data")

	setIfMissing(payload, "purchase_order_number", asString(payload["request_id"]))
	setIfMissing(payload, "reference_id", asString(payload["request_id"]))
	setIfMissing(payload, "issue_date", normalizeDateLike(asString(payload["created_at"])))
	setIfMissing(payload, "status_label", humanizeEnum(firstNonEmpty(payload["status"], serviceRaw["status"]), "SERVICE_REQUEST_STATUS_"))
	setIfMissing(payload, "request_type_label", humanizeEnum(firstNonEmpty(payload["request_type"], serviceRaw["request_type"]), "SERVICE_REQUEST_TYPE_"))
	setIfMissing(payload, "policy_reference", asString(payload["policy_id"]))
	setIfMissing(payload, "customer_reference", asString(payload["customer_id"]))
	setIfMissing(payload, "processed_date", normalizeDateLike(asString(payload["processed_at"])))
	setIfMissing(payload, "processed_by", asString(payload["processed_by"]))
	setIfMissing(payload, "customer_name", firstNonEmpty(payload["customer_name"], payload["requester_name"]))
	setIfMissing(payload, "customer_email", firstNonEmpty(payload["customer_email"], payload["requester_email"]))
	setIfMissing(payload, "customer_phone", firstNonEmpty(payload["customer_phone"], payload["requester_phone"]))

	processingFee := firstNonEmpty(payload["processing_fee"], payload["service_fee"])
	if processingFee != "" {
		setIfMissing(payload, "processing_fee", processingFee)
		setIfMissing(payload, "service_fee", processingFee)
	}

	if _, ok := payload["items"]; !ok {
		amount := valueOr(payload, "service_fee", valueOr(payload, "processing_fee", "0.00"))
		payload["items"] = []any{
			map[string]any{
				"description": valueOr(payload, "request_type_label", "Policy service request"),
				"quantity":    "1",
				"unit_price":  amount,
				"amount":      amount,
			},
		}
	}
}

func flattenEnvelope(payload map[string]any, key string) {
	raw, ok := payload[key]
	if !ok {
		return
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return
	}
	for k, v := range m {
		if _, exists := payload[k]; !exists {
			payload[k] = v
		}
	}
}

func mergeExcelTemplateData(templateName string, payload map[string]any) error {
	encoded := strings.TrimSpace(asString(payload[excelTemplateBase64Key]))
	if encoded == "" {
		return nil
	}

	rawBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("%w: invalid excel_template_base64 payload", ErrInvalidInput)
	}

	book, err := excelize.OpenReader(bytes.NewReader(rawBytes))
	if err != nil {
		return fmt.Errorf("%w: failed to open excel template", ErrInvalidInput)
	}
	defer func() { _ = book.Close() }()

	dataSheet := strings.TrimSpace(asString(payload[excelSheetKey]))
	if dataSheet == "" {
		sheets := book.GetSheetList()
		if len(sheets) == 0 {
			return fmt.Errorf("%w: excel workbook has no sheets", ErrInvalidInput)
		}
		dataSheet = sheets[0]
	}

	headerRow := intFromAny(payload[excelHeaderRowKey], 1)
	if headerRow < 1 {
		headerRow = 1
	}

	items, subtotal, err := parseLineItemsFromSheet(book, dataSheet, headerRow)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if len(items) > 0 {
		payload["items"] = items
		payload["subtotal"] = formatMoney(subtotal)
	}

	metaSheet := strings.TrimSpace(asString(payload[excelMetaSheetKey]))
	if metaSheet == "" {
		metaSheet = pickMetaSheet(book)
	}
	if metaSheet != "" {
		meta, metaErr := parseMetaSheet(book, metaSheet)
		if metaErr == nil {
			for k, v := range meta {
				if _, exists := payload[k]; !exists {
					payload[k] = v
				}
			}
		}
	}

	delete(payload, excelTemplateBase64Key)
	delete(payload, excelSheetKey)
	delete(payload, excelHeaderRowKey)
	delete(payload, excelMetaSheetKey)

	return nil
}

func parseLineItemsFromSheet(book *excelize.File, sheetName string, headerRow int) ([]any, float64, error) {
	rows, err := book.GetRows(sheetName)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read sheet %q", sheetName)
	}
	if len(rows) < headerRow {
		return nil, 0, fmt.Errorf("sheet %q has no header row %d", sheetName, headerRow)
	}

	header := rows[headerRow-1]
	if len(header) == 0 {
		return nil, 0, fmt.Errorf("sheet %q header is empty", sheetName)
	}

	colIndex := map[string]int{
		"description": findHeaderColumn(header, excelHeaderAliases["description"]...),
		"quantity":    findHeaderColumn(header, excelHeaderAliases["quantity"]...),
		"unit_price":  findHeaderColumn(header, excelHeaderAliases["unit_price"]...),
		"amount":      findHeaderColumn(header, excelHeaderAliases["amount"]...),
	}
	if colIndex["description"] < 0 {
		return nil, 0, fmt.Errorf("sheet %q missing description column", sheetName)
	}

	items := make([]any, 0)
	subtotal := 0.0
	for i := headerRow; i < len(rows); i++ {
		row := rows[i]
		description := cellAt(row, colIndex["description"])
		qtyRaw := cellAt(row, colIndex["quantity"])
		unitPriceRaw := cellAt(row, colIndex["unit_price"])
		amountRaw := cellAt(row, colIndex["amount"])
		if strings.TrimSpace(description) == "" &&
			strings.TrimSpace(qtyRaw) == "" &&
			strings.TrimSpace(unitPriceRaw) == "" &&
			strings.TrimSpace(amountRaw) == "" {
			continue
		}

		qty := numberFromString(qtyRaw, 1)
		unitPrice := numberFromString(unitPriceRaw, 0)
		amount := numberFromString(amountRaw, -1)
		if amount < 0 {
			amount = qty * unitPrice
		}

		items = append(items, map[string]any{
			"description": description,
			"quantity":    formatQuantity(qty),
			"unit_price":  formatMoney(unitPrice),
			"amount":      formatMoney(amount),
		})
		subtotal += amount
	}

	return items, subtotal, nil
}

func parseMetaSheet(book *excelize.File, sheetName string) (map[string]string, error) {
	rows, err := book.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	meta := make(map[string]string)
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		key := normalizeHeader(row[0])
		val := strings.TrimSpace(row[1])
		if key == "" || val == "" {
			continue
		}
		meta[key] = val
	}
	return meta, nil
}

func pickMetaSheet(book *excelize.File) string {
	for _, name := range book.GetSheetList() {
		lower := strings.ToLower(strings.TrimSpace(name))
		if lower == "meta" || lower == "metadata" || lower == "summary" {
			return name
		}
	}
	return ""
}

func normalizeTemplateTotals(templateName string, payload map[string]any) {
	kind := templateKind(templateName)
	subtotal := subtotalFromPayload(payload)
	tax := amountFromPayload(payload, "tax")
	if tax == 0 {
		rate := numberFromAny(payload[taxRateKey], 0)
		if rate > 0 {
			tax = subtotal * rate
		}
	}

	switch kind {
	case "invoice":
		total := amountFromPayload(payload, "total")
		if total == 0 {
			total = subtotal + tax
		}
		payload["subtotal"] = formatMoney(subtotal)
		payload["tax"] = formatMoney(tax)
		payload["total"] = formatMoney(total)
	case "purchase_order":
		serviceFee := amountFromPayload(payload, "service_fee")
		if serviceFee == 0 {
			serviceFee = amountFromPayload(payload, "processing_fee")
		}
		total := amountFromPayload(payload, "total")
		if total == 0 {
			total = subtotal + tax + serviceFee
		}
		payload["subtotal"] = formatMoney(subtotal)
		payload["service_fee"] = formatMoney(serviceFee)
		payload["processing_fee"] = formatMoney(serviceFee)
		payload["tax"] = formatMoney(tax)
		payload["total"] = formatMoney(total)
	}
}

func subtotalFromPayload(payload map[string]any) float64 {
	items, ok := payload["items"]
	if ok {
		if subtotal := subtotalFromItems(items); subtotal > 0 {
			return subtotal
		}
	}
	return amountFromPayload(payload, "subtotal")
}

func subtotalFromItems(items any) float64 {
	list, ok := items.([]any)
	if !ok {
		return 0
	}
	sum := 0.0
	for _, item := range list {
		entry, ok := item.(map[string]any)
		if !ok {
			continue
		}
		amount := amountFromAny(entry["amount"], -1)
		if amount < 0 {
			qty := numberFromAny(entry["quantity"], 1)
			unitPrice := numberFromAny(entry["unit_price"], 0)
			amount = qty * unitPrice
		}
		sum += amount
	}
	return sum
}

func amountFromPayload(payload map[string]any, key string) float64 {
	return amountFromAny(payload[key], 0)
}

func amountFromAny(v any, fallback float64) float64 {
	switch t := v.(type) {
	case nil:
		return fallback
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case string:
		return numberFromString(t, fallback)
	default:
		return numberFromString(asString(v), fallback)
	}
}

func numberFromAny(v any, fallback float64) float64 {
	return amountFromAny(v, fallback)
}

func intFromAny(v any, fallback int) int {
	switch t := v.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	case float32:
		return int(t)
	case float64:
		return int(t)
	case string:
		if t == "" {
			return fallback
		}
		n, err := strconv.Atoi(strings.TrimSpace(t))
		if err != nil {
			return fallback
		}
		return n
	default:
		return fallback
	}
}

func numberFromString(raw string, fallback float64) float64 {
	s := strings.TrimSpace(raw)
	if s == "" {
		return fallback
	}
	normalized := strings.ReplaceAll(s, ",", "")
	normalized = strings.ReplaceAll(normalized, "৳", "")
	normalized = strings.ReplaceAll(normalized, "$", "")
	val, err := strconv.ParseFloat(strings.TrimSpace(normalized), 64)
	if err != nil {
		return fallback
	}
	return val
}

func billingAmountDecimal(v any) (float64, bool) {
	switch t := v.(type) {
	case nil:
		return 0, false
	case map[string]any:
		decimalVal := amountFromAny(t["decimal_amount"], 0)
		if decimalVal > 0 {
			return decimalVal, true
		}
		minorVal := amountFromAny(t["amount"], 0)
		if minorVal > 0 {
			return minorVal / 100.0, true
		}
		return 0, false
	case float64, float32, int, int32, int64:
		return amountFromAny(t, 0), true
	case string:
		parsed := numberFromString(t, 0)
		if parsed == 0 {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func normalizeDateLike(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "T") && len(s) >= 10 {
		return s[:10]
	}
	if len(s) >= 10 && s[4] == '-' && s[7] == '-' {
		return s[:10]
	}
	return s
}

func moneyToDisplay(v any) string {
	if v == nil {
		return ""
	}
	if m, ok := v.(map[string]any); ok {
		if dec := amountFromAny(m["decimal_amount"], 0); dec > 0 {
			return formatMoney(dec)
		}
		if minor := amountFromAny(m["amount"], 0); minor > 0 {
			return formatMoney(minor / 100.0)
		}
	}
	if parsed := amountFromAny(v, 0); parsed > 0 {
		return formatMoney(parsed)
	}
	return ""
}

func valueOr(payload map[string]any, key, fallback string) string {
	if payload == nil {
		return fallback
	}
	if v := strings.TrimSpace(asString(payload[key])); v != "" {
		return v
	}
	return fallback
}

func setIfMissingString(payload map[string]any, key, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	setIfMissing(payload, key, value)
}

func firstNonEmpty(values ...any) string {
	for _, v := range values {
		s := strings.TrimSpace(asString(v))
		if s != "" {
			return s
		}
	}
	return ""
}

func humanizeEnum(raw string, prefixes ...string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	normalized := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_"))
	for _, prefix := range prefixes {
		p := strings.ToUpper(strings.TrimSpace(prefix))
		if p == "" {
			continue
		}
		if strings.HasPrefix(normalized, p) {
			normalized = strings.TrimPrefix(normalized, p)
			break
		}
	}
	parts := strings.Split(normalized, "_")
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		lower := strings.ToLower(part)
		words = append(words, strings.ToUpper(lower[:1])+lower[1:])
	}
	if len(words) == 0 {
		return strings.TrimSpace(raw)
	}
	return strings.Join(words, " ")
}

func toStringSlice(v any) []string {
	switch t := v.(type) {
	case nil:
		return nil
	case []string:
		return t
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			s := strings.TrimSpace(asString(item))
			if s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		s := strings.TrimSpace(asString(v))
		if s == "" {
			return nil
		}
		return []string{s}
	}
}

func joinStringSlice(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, sep)
}

func formatMoney(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func formatQuantity(v float64) string {
	if math.Abs(v-math.Round(v)) < 1e-9 {
		return strconv.FormatInt(int64(math.Round(v)), 10)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func findHeaderColumn(header []string, aliases ...string) int {
	aliasSet := make(map[string]struct{}, len(aliases))
	for _, a := range aliases {
		aliasSet[normalizeHeader(a)] = struct{}{}
	}
	for idx, raw := range header {
		if _, ok := aliasSet[normalizeHeader(raw)]; ok {
			return idx
		}
	}
	return -1
}

func normalizeHeader(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func cellAt(row []string, col int) string {
	if col < 0 || col >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[col])
}
