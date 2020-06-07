package main

import (
	"fmt"
	"google.golang.org/api/monitoring/v3"
	"strings"
)

type Result struct {
	Header []string
	Rows   []Row
}

type Row struct {
	Columns []string
}

func buildQueryResult(descriptor *monitoring.TimeSeriesDescriptor, history []*monitoring.TimeSeriesData) *Result {
	result := Result{}

	// Header columns for label
	for _, ld := range descriptor.LabelDescriptors {
		result.Header = append(result.Header, normalizeLabelKey(ld.Key))
	}

	// Header columns for time
	var showStartTime bool
	for _, pd := range descriptor.PointDescriptors {
		if pd.MetricKind == "DELTA" || pd.MetricKind == "CUMULATIVE" {
			showStartTime = true
		}
	}
	if showStartTime {
		result.Header = append(result.Header, "start_time", "end_time")
	} else {
		result.Header = append(result.Header, "time")
	}

	// Header columns for value
	for _, pd := range descriptor.PointDescriptors {
		result.Header = append(result.Header, normalizeValueKey(pd.Key))
	}

	// Body columns
	for _, data := range history {
		var labelValues [] string
		for i, lv := range data.LabelValues {
			valueType := descriptor.LabelDescriptors[i].ValueType
			encoded := encodeLabelValue(valueType, lv)
			labelValues = append(labelValues, encoded)
		}

		for _, point := range data.PointData {
			var row Row
			row.Columns = append(row.Columns, labelValues...)

			if showStartTime {
				row.Columns = append(row.Columns, point.TimeInterval.StartTime, point.TimeInterval.EndTime)
			} else {
				row.Columns = append(row.Columns, point.TimeInterval.EndTime)
			}

			for i, v := range point.Values {
				valueType := descriptor.PointDescriptors[i].ValueType
				encoded := encodeValue(valueType, v)
				row.Columns = append(row.Columns, encoded)
			}
			result.Rows = append(result.Rows, row)
		}
	}
	return &result
}

func normalizeLabelKey(key string) string {
	// remove prefix for monitored resource
	key = strings.TrimPrefix(key, "resource.")
	// remove prefix for metric
	key = strings.TrimPrefix(key, "metric.")
	return key
}

func normalizeValueKey(key string) string {
	return strings.TrimPrefix(key, "value.")
}

func encodeLabelValue(valueType string, value *monitoring.LabelValue) string {
	switch valueType {
	case "STRING", "": // Empty valueType is STRING by default
		return value.StringValue
	case "BOOL":
		return fmt.Sprintf("%t", value.BoolValue)
	case "INT64":
		return fmt.Sprintf("%d", value.Int64Value)
	default:
		return ""
	}
}

func encodeValue(valueType string, value *monitoring.TypedValue) string {
	switch valueType {
	case "VALUE_TYPE_UNSPECIFIED":
		return ""
	case "BOOL":
		if value.BoolValue == nil {
			return ""
		}
		return fmt.Sprintf("%t", *value.BoolValue)
	case "INT64":
		if value.Int64Value == nil {
			return ""
		}
		return fmt.Sprintf("%d", *value.Int64Value)
	case "DOUBLE":
		if value.DoubleValue == nil {
			return ""
		}
		return fmt.Sprintf("%f", *value.DoubleValue)
	case "STRING":
		if value.StringValue == nil {
			return ""
		}
		return *value.StringValue
	case "DISTRIBUTION":
		if value.DistributionValue == nil {
			return ""
		}
		return fmt.Sprintf("%f", value.DistributionValue.Mean)
	default:
		return ""
	}
}
