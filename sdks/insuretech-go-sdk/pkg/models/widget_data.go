package models


// WidgetData represents a widget_data
type WidgetData struct {
	Rows []*Row `json:"rows,omitempty"`
	Columns []string `json:"columns,omitempty"`
}
