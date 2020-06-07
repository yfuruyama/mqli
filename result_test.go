package main

import (
	"google.golang.org/api/monitoring/v3"
	"testing"
)

func TestEncodeLabelValue(t *testing.T) {
	for _, tt := range []struct {
		desc      string
		valueType string
		value     *monitoring.LabelValue
		want      string
	}{
		{
			desc:      "valueType is empty",
			valueType: "",
			value: &monitoring.LabelValue{
				StringValue: "foo",
			},
			want: "foo",
		},
		{
			desc:      "string",
			valueType: "STRING",
			value: &monitoring.LabelValue{
				StringValue: "bar",
			},
			want: "bar",
		},
		{
			desc:      "bool",
			valueType: "BOOL",
			value: &monitoring.LabelValue{
				BoolValue: true,
			},
			want: "true",
		},
		{
			desc:      "int64",
			valueType: "INT64",
			value: &monitoring.LabelValue{
				Int64Value: 123,
			},
			want: "123",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			if got := encodeLabelValue(tt.valueType, tt.value); got != tt.want {
				t.Errorf("encodeLabelValue(%q, %v) = %q, but want %q", tt.valueType, tt.value, got, tt.want)
			}
		})
	}
}

func TestEncodeValue(t *testing.T) {
	for _, tt := range []struct {
		desc      string
		valueType string
		value     *monitoring.TypedValue
		want      string
	}{
		{
			desc:      "empty bool",
			valueType: "BOOL",
			value: &monitoring.TypedValue{
				BoolValue: nil,
			},
			want: "",
		},
		{
			desc:      "bool",
			valueType: "BOOL",
			value: &monitoring.TypedValue{
				BoolValue: boolPointer(true),
			},
			want: "true",
		},
		{
			desc:      "empty int64",
			valueType: "INT64",
			value: &monitoring.TypedValue{
				Int64Value: nil,
			},
			want: "",
		},
		{
			desc:      "int64",
			valueType: "INT64",
			value: &monitoring.TypedValue{
				Int64Value: int64Pointer(123),
			},
			want: "123",
		},
		{
			desc:      "empty double",
			valueType: "DOUBLE",
			value: &monitoring.TypedValue{
				DoubleValue: nil,
			},
			want: "",
		},
		{
			desc:      "double",
			valueType: "DOUBLE",
			value: &monitoring.TypedValue{
				DoubleValue: float64Pointer(1.23),
			},
			want: "1.230000",
		},
		{
			desc:      "empty string",
			valueType: "STRING",
			value: &monitoring.TypedValue{
				StringValue: nil,
			},
			want: "",
		},
		{
			desc:      "string",
			valueType: "STRING",
			value: &monitoring.TypedValue{
				StringValue: stringPointer("foo"),
			},
			want: "foo",
		},
		{
			desc:      "empty distribution",
			valueType: "DISTRIBUTION",
			value: &monitoring.TypedValue{
				DistributionValue: nil,
			},
			want: "",
		},
		{
			desc:      "distribution",
			valueType: "DISTRIBUTION",
			value: &monitoring.TypedValue{
				DistributionValue: &monitoring.Distribution{
					Mean: 1.23,
				},
			},
			want: "1.230000",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			if got := encodeValue(tt.valueType, tt.value); got != tt.want {
				t.Errorf("encodeValue(%q, %v) = %q, but want %q", tt.valueType, tt.value, got, tt.want)
			}
		})
	}
}

func boolPointer(v bool) *bool {
	return &v
}

func int64Pointer(v int64) *int64 {
	return &v
}

func float64Pointer(v float64) *float64 {
	return &v
}

func stringPointer(v string) *string {
	return &v
}
