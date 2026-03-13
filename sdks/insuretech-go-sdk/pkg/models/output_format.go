package models

// OutputFormat represents a output_format
type OutputFormat string

// OutputFormat values
const (
	OutputFormatOUTPUTFORMATUNSPECIFIED OutputFormat = "OUTPUT_FORMAT_UNSPECIFIED"
	OutputFormatOUTPUTFORMATPDF  = "OUTPUT_FORMAT_PDF"
	OutputFormatOUTPUTFORMATHTML  = "OUTPUT_FORMAT_HTML"
	OutputFormatOUTPUTFORMATDOCX  = "OUTPUT_FORMAT_DOCX"
)
