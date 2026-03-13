package models


// IndexOptions represents a index_options
type IndexOptions struct {
	IndexType *IndexType `json:"index_type,omitempty"`
	Unique bool `json:"unique,omitempty"`
	CompositeFields []string `json:"composite_fields,omitempty"`
	IndexMethod string `json:"index_method,omitempty"`
	WhereClause string `json:"where_clause,omitempty"`
	IndexName string `json:"index_name,omitempty"`
}
