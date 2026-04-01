package models


// WidgetData represents a widget_data
type WidgetData struct {
	Columns []string `json:"columns,omitempty"`
	Rows []*Row `json:"rows,omitempty"`
}
