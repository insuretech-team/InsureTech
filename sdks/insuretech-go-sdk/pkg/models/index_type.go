package models

// IndexType represents a index_type
type IndexType string

// IndexType values
const (
	IndexTypeINDEXTYPEUNSPECIFIED IndexType = "INDEX_TYPE_UNSPECIFIED"
	IndexTypeINDEXTYPENONE  = "INDEX_TYPE_NONE"
	IndexTypeINDEXTYPEBTREE  = "INDEX_TYPE_BTREE"
	IndexTypeINDEXTYPEHASH  = "INDEX_TYPE_HASH"
	IndexTypeINDEXTYPEGIN  = "INDEX_TYPE_GIN"
	IndexTypeINDEXTYPEGIST  = "INDEX_TYPE_GIST"
	IndexTypeINDEXTYPEBRIN  = "INDEX_TYPE_BRIN"
)
