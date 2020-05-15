package main

import (
	"context"
	"google.golang.org/api/monitoring/v3"
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

	result := buildQueryResult(resp.TimeSeriesDescriptor, resp.TimeSeriesData)
	return result, nil
}
