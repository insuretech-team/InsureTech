package service

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestEnrichTemplatePayloadInvoiceFromExcel(t *testing.T) {
	payload := map[string]any{
		excelTemplateBase64Key: buildWorkbookBase64(t,
			[][]string{
				{"Description", "Qty", "Unit Price"},
				{"Policy Premium", "2", "100"},
				{"Service Fee", "1", "20"},
			},
			map[string]string{
				"invoice_number": "INV-20260305-001",
				"customer_name":  "Demo Customer",
			},
		),
		"tax_rate": 0.10,
	}

	err := enrichTemplatePayload("invoice", payload)
	require.NoError(t, err)

	items, ok := payload["items"].([]any)
	require.True(t, ok)
	require.Len(t, items, 2)

	require.Equal(t, "220.00", payload["subtotal"])
	require.Equal(t, "22.00", payload["tax"])
	require.Equal(t, "242.00", payload["total"])
	require.Equal(t, "INV-20260305-001", payload["invoice_number"])
	require.Equal(t, "Demo Customer", payload["customer_name"])

	_, hasExcelBlob := payload[excelTemplateBase64Key]
	require.False(t, hasExcelBlob)
}

func TestEnrichTemplatePayloadPurchaseOrderFromExcel(t *testing.T) {
	payload := map[string]any{
		excelTemplateBase64Key: buildWorkbookBase64(t,
			[][]string{
				{"Item", "Quantity", "Unit Cost"},
				{"Printer Paper", "3", "50"},
				{"Stapler", "1", "10"},
			},
			map[string]string{
				"purchase_order_number": "PO-20260305-001",
				"vendor_name":           "Stationery House",
			},
		),
		"service_fee": "15",
		"tax_rate":    "0.05",
	}

	err := enrichTemplatePayload("purchase_order", payload)
	require.NoError(t, err)

	require.Equal(t, "160.00", payload["subtotal"])
	require.Equal(t, "8.00", payload["tax"])
	require.Equal(t, "15.00", payload["service_fee"])
	require.Equal(t, "15.00", payload["processing_fee"])
	require.Equal(t, "183.00", payload["total"])
	require.Equal(t, "PO-20260305-001", payload["purchase_order_number"])
	require.Equal(t, "Stationery House", payload["vendor_name"])
}

func TestNormalizeTemplateTotalsFromInlineItems(t *testing.T) {
	payload := map[string]any{
		"items": []any{
			map[string]any{
				"description": "Line A",
				"quantity":    "2",
				"unit_price":  "50",
			},
			map[string]any{
				"description": "Line B",
				"amount":      "25",
			},
		},
	}

	normalizeTemplateTotals("invoice", payload)

	require.Equal(t, "125.00", payload["subtotal"])
	require.Equal(t, "0.00", payload["tax"])
	require.Equal(t, "125.00", payload["total"])
}

func buildWorkbookBase64(t *testing.T, itemRows [][]string, meta map[string]string) string {
	t.Helper()

	wb := excelize.NewFile()
	itemsSheet := "Items"
	wb.SetSheetName("Sheet1", itemsSheet)

	for r := range itemRows {
		for c := range itemRows[r] {
			cell, err := excelize.CoordinatesToCellName(c+1, r+1)
			require.NoError(t, err)
			err = wb.SetCellValue(itemsSheet, cell, itemRows[r][c])
			require.NoError(t, err)
		}
	}

	metaSheet := "Meta"
	_, err := wb.NewSheet(metaSheet)
	require.NoError(t, err)
	row := 1
	for k, v := range meta {
		keyCell, keyErr := excelize.CoordinatesToCellName(1, row)
		require.NoError(t, keyErr)
		valCell, valErr := excelize.CoordinatesToCellName(2, row)
		require.NoError(t, valErr)
		require.NoError(t, wb.SetCellValue(metaSheet, keyCell, k))
		require.NoError(t, wb.SetCellValue(metaSheet, valCell, v))
		row++
	}

	buf, err := wb.WriteToBuffer()
	require.NoError(t, err)
	require.NoError(t, wb.Close())

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
