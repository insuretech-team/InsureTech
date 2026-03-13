package models

// MetricType represents a metric_type
type MetricType string

// MetricType values
const (
	MetricTypeMETRICTYPEUNSPECIFIED MetricType = "METRIC_TYPE_UNSPECIFIED"
	MetricTypeMETRICTYPECOUNTER  = "METRIC_TYPE_COUNTER"
	MetricTypeMETRICTYPEGAUGE  = "METRIC_TYPE_GAUGE"
	MetricTypeMETRICTYPEHISTOGRAM  = "METRIC_TYPE_HISTOGRAM"
)
