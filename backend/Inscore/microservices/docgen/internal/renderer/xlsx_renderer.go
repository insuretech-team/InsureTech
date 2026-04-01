// Package renderer provides rich document output renderers for the docgen service.
// xlsx_renderer.go — Production-quality Excel (.xlsx) generation using excelize.
//
// Features:
//   - Branded header row (bold, white text, deep-blue background)
//   - Alternating row shading for readability
//   - Currency number format for money columns
//   - Auto-fitted column widths (capped at 60 chars)
//   - Frozen header row (freeze pane)
//   - Multi-sheet output: Items, Summary, Meta
//   - Print area + A4 page setup with print titles
//   - Document metadata (title, author, created date)
//   - Right-to-left number alignment in numeric cells
package renderer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// XLSXOptions controls how the spreadsheet is built.
type XLSXOptions struct {
	Title       string
	Author      string
	Subject     string
	Description string
	// Items is the primary line-item list ([]map[string]any).
	Items []map[string]any
	// ItemColumns defines column order and headers for the Items sheet.
	// Each entry: {Key: "description", Header: "Description", Width: 40, IsMoney: false}
	ItemColumns []XLSXColumn
	// Summary is a flat key-value map written to the Summary sheet.
	Summary map[string]any
	// Meta is an optional second key-value map for document metadata.
	Meta map[string]string
	// Totals are written at the bottom of the Items sheet.
	// Each entry: {Label: "Total", Value: "1,234.56", IsBold: true}
	Totals []XLSXTotalRow
	// SheetName overrides the default "Items" sheet name.
	SheetName string
}

// XLSXColumn describes one column in the Items sheet.
type XLSXColumn struct {
	Key     string  // map key in item
	Header  string  // displayed column header
	Width   float64 // column width in characters (0 = auto)
	IsMoney bool    // apply currency format
	IsDate  bool    // apply date format
}

// XLSXTotalRow is one row in the totals block.
type XLSXTotalRow struct {
	Label  string
	Value  string
	IsBold bool
}

// ─── Style constants ─────────────────────────────────────────────────────────

const (
	colorHeaderBG   = "0D47A1" // deep brand blue
	colorHeaderFG   = "FFFFFF"
	colorAltRow     = "EEF2F7" // light blue-grey
	colorTotalsBG   = "E3EAF3"
	colorFinalBG    = "0D47A1"
	colorFinalFG    = "FFFFFF"
	colorBorderLine = "B0C4DE"
	colorSummaryKey = "1565C0"
)

// RenderXLSX generates a styled Excel workbook and returns the raw bytes.
func RenderXLSX(opts XLSXOptions) ([]byte, error) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	// ── Document properties ───────────────────────────────────────────────────
	_ = f.SetDocProps(&excelize.DocProperties{
		Title:       opts.Title,
		Subject:     opts.Subject,
		Creator:     opts.Author,
		Description: opts.Description,
		Created:     time.Now().UTC().Format(time.RFC3339),
		Modified:    time.Now().UTC().Format(time.RFC3339),
	})

	// ── Pre-build shared styles ───────────────────────────────────────────────
	styles, err := buildStyles(f)
	if err != nil {
		return nil, fmt.Errorf("xlsx: failed to build styles: %w", err)
	}

	// ── Items sheet ───────────────────────────────────────────────────────────
	sheetName := opts.SheetName
	if strings.TrimSpace(sheetName) == "" {
		sheetName = "Items"
	}
	// excelize creates Sheet1 by default; rename it.
	f.SetSheetName("Sheet1", sheetName)

	if err := writeItemsSheet(f, sheetName, opts, styles); err != nil {
		return nil, err
	}

	// ── Summary sheet ─────────────────────────────────────────────────────────
	if len(opts.Summary) > 0 || len(opts.Totals) > 0 {
		summarySheet := "Summary"
		_, _ = f.NewSheet(summarySheet)
		if err := writeSummarySheet(f, summarySheet, opts, styles); err != nil {
			return nil, err
		}
	}

	// ── Meta sheet ────────────────────────────────────────────────────────────
	if len(opts.Meta) > 0 {
		metaSheet := "Meta"
		_, _ = f.NewSheet(metaSheet)
		if err := writeMetaSheet(f, metaSheet, opts.Meta, styles); err != nil {
			return nil, err
		}
	}

	// ── Set active sheet to Items ─────────────────────────────────────────────
	idx, _ := f.GetSheetIndex(sheetName)
	f.SetActiveSheet(idx)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("xlsx: failed to write buffer: %w", err)
	}
	return buf.Bytes(), nil
}

// ─── Style registry ──────────────────────────────────────────────────────────

type styleIDs struct {
	header      int
	normal      int
	alt         int
	money       int
	moneyAlt    int
	bold        int
	totalsLabel int
	totalsValue int
	finalLabel  int
	finalValue  int
	summaryKey  int
	summaryVal  int
	dateStyle   int
}

func buildStyles(f *excelize.File) (*styleIDs, error) {
	s := &styleIDs{}
	var err error

	border := []excelize.Border{
		{Type: "left",   Color: colorBorderLine, Style: 1},
		{Type: "right",  Color: colorBorderLine, Style: 1},
		{Type: "top",    Color: colorBorderLine, Style: 1},
		{Type: "bottom", Color: colorBorderLine, Style: 1},
	}
	thickBorderBottom := []excelize.Border{
		{Type: "left",   Color: colorBorderLine, Style: 1},
		{Type: "right",  Color: colorBorderLine, Style: 1},
		{Type: "top",    Color: colorBorderLine, Style: 1},
		{Type: "bottom", Color: colorHeaderBG,   Style: 2},
	}

	// Header style
	s.header, err = f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Bold: true, Color: colorHeaderFG, Size: 11},
		Fill:   excelize.Fill{Type: "pattern", Color: []string{colorHeaderBG}, Pattern: 1},
		Border: thickBorderBottom,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return nil, err
	}

	// Normal row
	s.normal, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Border:    border,
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Alt row
	s.alt, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorAltRow}, Pattern: 1},
		Border:    border,
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Money normal
	s.money, err = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 10},
		Border:       border,
		NumFmt:       4, // #,##0.00
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Money alt
	s.moneyAlt, err = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 10},
		Fill:         excelize.Fill{Type: "pattern", Color: []string{colorAltRow}, Pattern: 1},
		Border:       border,
		NumFmt:       4,
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Bold
	s.bold, err = f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Bold: true, Size: 10},
		Border: border,
	})
	if err != nil {
		return nil, err
	}

	// Totals label
	s.totalsLabel, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: colorSummaryKey},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorTotalsBG}, Pattern: 1},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Totals value
	s.totalsValue, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorTotalsBG}, Pattern: 1},
		Border:    border,
		NumFmt:    4,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Final total label (dark bg)
	s.finalLabel, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: colorHeaderFG},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorFinalBG}, Pattern: 1},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Final total value
	s.finalValue, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: colorHeaderFG},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorFinalBG}, Pattern: 1},
		Border:    border,
		NumFmt:    4,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Summary key
	s.summaryKey, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: colorSummaryKey},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorTotalsBG}, Pattern: 1},
		Border:    border,
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Summary value
	s.summaryVal, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Border:    border,
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Date
	s.dateStyle, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Border:    border,
		NumFmt:    14, // m/d/yy
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

// ─── Items sheet ──────────────────────────────────────────────────────────────

func writeItemsSheet(f *excelize.File, sheet string, opts XLSXOptions, styles *styleIDs) error {
	cols := opts.ItemColumns
	if len(cols) == 0 {
		// Auto-detect columns from first item
		if len(opts.Items) > 0 {
			for k := range opts.Items[0] {
				cols = append(cols, XLSXColumn{Key: k, Header: humanizeKey(k)})
			}
		}
	}
	if len(cols) == 0 {
		return nil
	}

	// ── Header row ────────────────────────────────────────────────────────────
	for colIdx, col := range cols {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		_ = f.SetCellValue(sheet, cell, col.Header)
		_ = f.SetCellStyle(sheet, cell, cell, styles.header)
	}

	// ── Data rows ─────────────────────────────────────────────────────────────
	maxWidths := make([]float64, len(cols))
	for colIdx, col := range cols {
		w := float64(len(col.Header))
		if col.Width > 0 {
			w = col.Width
		}
		maxWidths[colIdx] = w
	}

	dataRowStart := 2
	for rowIdx, item := range opts.Items {
		excelRow := rowIdx + dataRowStart
		isAlt := rowIdx%2 == 1
		for colIdx, col := range cols {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, excelRow)
			val := item[col.Key]
			styleID := styles.normal
			if isAlt {
				styleID = styles.alt
			}
			if col.IsMoney {
				styleID = styles.money
				if isAlt {
					styleID = styles.moneyAlt
				}
				num := parseFloat(fmt.Sprintf("%v", val))
				_ = f.SetCellValue(sheet, cell, num)
			} else if col.IsDate {
				styleID = styles.dateStyle
				_ = f.SetCellValue(sheet, cell, fmt.Sprintf("%v", val))
			} else {
				_ = f.SetCellValue(sheet, cell, fmt.Sprintf("%v", val))
			}
			_ = f.SetCellStyle(sheet, cell, cell, styleID)

			// Track max width for auto-fit
			valStr := fmt.Sprintf("%v", val)
			if float64(len(valStr)) > maxWidths[colIdx] {
				maxWidths[colIdx] = math.Min(float64(len(valStr)), 60)
			}
		}
	}

	// ── Totals block ──────────────────────────────────────────────────────────
	totalsStartRow := dataRowStart + len(opts.Items) + 1
	if len(opts.Totals) > 0 {
		// Merge left columns to create label space
		labelEndCol := len(cols) - 1
		if labelEndCol < 1 {
			labelEndCol = 1
		}
		for tIdx, total := range opts.Totals {
			excelRow := totalsStartRow + tIdx
			// Label cell (merged across all but last column)
			labelCell, _ := excelize.CoordinatesToCellName(1, excelRow)
			labelEndCell, _ := excelize.CoordinatesToCellName(labelEndCol, excelRow)
			valueCell, _ := excelize.CoordinatesToCellName(len(cols), excelRow)

			if labelCell != labelEndCell {
				_ = f.MergeCell(sheet, labelCell, labelEndCell)
			}

			labelStyleID := styles.totalsLabel
			valueStyleID := styles.totalsValue
			if total.IsBold {
				labelStyleID = styles.finalLabel
				valueStyleID = styles.finalValue
			}

			_ = f.SetCellValue(sheet, labelCell, total.Label)
			_ = f.SetCellStyle(sheet, labelCell, labelEndCell, labelStyleID)

			num := parseFloat(total.Value)
			_ = f.SetCellValue(sheet, valueCell, num)
			_ = f.SetCellStyle(sheet, valueCell, valueCell, valueStyleID)
		}
	}

	// ── Column widths ─────────────────────────────────────────────────────────
	for colIdx, col := range cols {
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		w := maxWidths[colIdx] + 2
		if col.Width > 0 {
			w = col.Width
		}
		_ = f.SetColWidth(sheet, colName, colName, w)
	}

	// ── Freeze header row ─────────────────────────────────────────────────────
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// ── Row height for header ─────────────────────────────────────────────────
	_ = f.SetRowHeight(sheet, 1, 22)

	// ── Page setup: A4, landscape for wide tables ─────────────────────────────
	landscape := len(cols) > 6
	pageSize := 9 // A4
	orientation := "portrait"
	if landscape {
		orientation = "landscape"
	}
	fitToHeight := 0
	fitToWidth := 1
	_ = f.SetPageLayout(sheet, &excelize.PageLayoutOptions{
		Size:        &pageSize,
		Orientation: &orientation,
		FitToHeight: &fitToHeight,
		FitToWidth:  &fitToWidth,
	})
	_ = f.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		OddHeader: "&C&B" + opts.Title,
		OddFooter: "&LInsureTech&CPage &P of &N&R" + time.Now().Format("2006-01-02"),
	})

	// Print titles (repeat header row on each printed page)
	_ = f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Titles",
		RefersTo: fmt.Sprintf("%s!$1:$1", sheet),
		Scope:    sheet,
	})

	// ── Print area ────────────────────────────────────────────────────────────
	lastDataRow := dataRowStart + len(opts.Items) + len(opts.Totals)
	lastCol, _ := excelize.ColumnNumberToName(len(cols))
	_ = f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: fmt.Sprintf("%s!$A$1:$%s$%d", sheet, lastCol, lastDataRow),
		Scope:    sheet,
	})

	return nil
}

// ─── Summary sheet ────────────────────────────────────────────────────────────

func writeSummarySheet(f *excelize.File, sheet string, opts XLSXOptions, styles *styleIDs) error {
	// Title
	_ = f.SetCellValue(sheet, "A1", opts.Title)
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: colorHeaderBG},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	_ = f.SetCellStyle(sheet, "A1", "A1", titleStyle)
	_ = f.MergeCell(sheet, "A1", "B1")
	_ = f.SetRowHeight(sheet, 1, 24)

	row := 2
	// Write summary key-value pairs
	for k, v := range opts.Summary {
		if strings.HasPrefix(k, "_") {
			continue
		}
		cellKey, _ := excelize.CoordinatesToCellName(1, row)
		cellVal, _ := excelize.CoordinatesToCellName(2, row)
		_ = f.SetCellValue(sheet, cellKey, humanizeKey(k))
		_ = f.SetCellValue(sheet, cellVal, fmt.Sprintf("%v", v))
		_ = f.SetCellStyle(sheet, cellKey, cellKey, styles.summaryKey)
		_ = f.SetCellStyle(sheet, cellVal, cellVal, styles.summaryVal)
		row++
	}

	// Blank row separator
	row++

	// Totals
	for _, total := range opts.Totals {
		cellKey, _ := excelize.CoordinatesToCellName(1, row)
		cellVal, _ := excelize.CoordinatesToCellName(2, row)
		_ = f.SetCellValue(sheet, cellKey, total.Label)
		num := parseFloat(total.Value)
		_ = f.SetCellValue(sheet, cellVal, num)
		labelStyleID := styles.summaryKey
		valueStyleID := styles.summaryVal
		if total.IsBold {
			labelStyleID = styles.finalLabel
			valueStyleID = styles.finalValue
		}
		_ = f.SetCellStyle(sheet, cellKey, cellKey, labelStyleID)
		_ = f.SetCellStyle(sheet, cellVal, cellVal, valueStyleID)
		row++
	}

	_ = f.SetColWidth(sheet, "A", "A", 28)
	_ = f.SetColWidth(sheet, "B", "B", 22)
	return nil
}

// ─── Meta sheet ───────────────────────────────────────────────────────────────

func writeMetaSheet(f *excelize.File, sheet string, meta map[string]string, styles *styleIDs) error {
	// Header
	_ = f.SetCellValue(sheet, "A1", "Key")
	_ = f.SetCellValue(sheet, "B1", "Value")
	_ = f.SetCellStyle(sheet, "A1", "B1", styles.header)
	_ = f.SetRowHeight(sheet, 1, 20)

	row := 2
	for k, v := range meta {
		cellKey, _ := excelize.CoordinatesToCellName(1, row)
		cellVal, _ := excelize.CoordinatesToCellName(2, row)
		_ = f.SetCellValue(sheet, cellKey, k)
		_ = f.SetCellValue(sheet, cellVal, v)
		styleID := styles.normal
		if row%2 == 0 {
			styleID = styles.alt
		}
		_ = f.SetCellStyle(sheet, cellKey, cellKey, styleID)
		_ = f.SetCellStyle(sheet, cellVal, cellVal, styleID)
		row++
	}

	_ = f.SetColWidth(sheet, "A", "A", 24)
	_ = f.SetColWidth(sheet, "B", "B", 40)
	return nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "৳", "")
	s = strings.ReplaceAll(s, "$", "")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func humanizeKey(k string) string {
	k = strings.ReplaceAll(k, "_", " ")
	words := strings.Fields(k)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, " ")
}
