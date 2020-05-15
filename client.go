package main

import (
	"context"
	"fmt"
	"google.golang.org/api/monitoring/v3"
	"strconv"
)

type Client struct {
	projectID string
}

func (c *Client) Query(q string) (*Result, error) {
	ctx := context.Background()
	s, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}
	svc := monitoring.NewProjectsTimeSeriesService(s)
	call := svc.Query("projects/"+c.projectID, &monitoring.QueryTimeSeriesRequest{
		Query: q,
	})
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	var labelKeys []string
	for _, ld := range resp.TimeSeriesDescriptor.LabelDescriptors {
		labelKeys = append(labelKeys, ld.Key)
	}
	var pointKeys []string
	for _, pd := range resp.TimeSeriesDescriptor.PointDescriptors {
		pointKeys = append(pointKeys, pd.Key)
	}
	var header []string
	header = append(header, labelKeys...)
	header = append(header, "start_time", "end_time")
	header = append(header, pointKeys...)

	result := Result{
		Header: header,
	}
	for _, data := range resp.TimeSeriesData {
		var labelValues [] string
		for _, lv := range data.LabelValues {
			labelValues = append(labelValues, lv.StringValue)
		}
		for _, point := range data.PointData {
			var row Row
			row.Columns = append(row.Columns, labelValues...)

			startTime := point.TimeInterval.StartTime
			endTime := point.TimeInterval.EndTime
			row.Columns = append(row.Columns, startTime, endTime)

			var values [] string
			for _, v := range point.Values {
				if v.StringValue != nil {
					values = append(values, *v.StringValue)
				}
				if v.DoubleValue != nil {
					values = append(values, strconv.FormatFloat(*v.DoubleValue, 'g', -1, 64))
				}
				if v.Int64Value != nil {
					values = append(values, fmt.Sprintf("%d", *v.Int64Value))
				}
			}
			row.Columns = append(row.Columns, values...)
			result.Rows = append(result.Rows, row)
		}
	}
	return &result, nil
}

type Result struct {
	Header []string
	Rows []Row
}

type Row struct {
	Columns []string
}
