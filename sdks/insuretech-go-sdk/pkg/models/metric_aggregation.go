package models

// MetricAggregation represents a metric_aggregation
type MetricAggregation string

// MetricAggregation values
const (
	MetricAggregationMETRICAGGREGATIONUNSPECIFIED MetricAggregation = "METRIC_AGGREGATION_UNSPECIFIED"
	MetricAggregationMETRICAGGREGATIONSUM  = "METRIC_AGGREGATION_SUM"
	MetricAggregationMETRICAGGREGATIONAVG  = "METRIC_AGGREGATION_AVG"
	MetricAggregationMETRICAGGREGATIONMIN  = "METRIC_AGGREGATION_MIN"
	MetricAggregationMETRICAGGREGATIONMAX  = "METRIC_AGGREGATION_MAX"
	MetricAggregationMETRICAGGREGATIONCOUNT  = "METRIC_AGGREGATION_COUNT"
)
