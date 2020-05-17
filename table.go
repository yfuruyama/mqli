package main

import (
	"fmt"
	"google.golang.org/api/monitoring/v3"
	"strconv"
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

	// Row columns
	for _, data := range history {
		var labelValues [] string
		for _, lv := range data.LabelValues {
			labelValues = append(labelValues, lv.StringValue)
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
				encoded := valueEncodeString(valueType, v)
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

func valueEncodeString(valueType string, value *monitoring.TypedValue) string {
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
		return strconv.FormatFloat(*value.DoubleValue, 'g', -1, 64)
	case "STRING":
		if value.StringValue == nil {
			return ""
		}
		return *value.StringValue
	case "DISTRIBUTION":
		if value.DistributionValue == nil {
			return ""
		}
		// TODO
		return "TODO"
	case "MONEY":
		// TODO
		return "TODO"
	default:
		return ""
	}
}
