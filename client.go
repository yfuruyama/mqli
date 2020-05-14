package main

import (
	"context"
	"google.golang.org/api/monitoring/v3"
)

type Client struct {
	projectID string
}

func (c *Client) Query(q string) (*monitoring.QueryTimeSeriesResponse, error) {
	ctx := context.Background()
	s, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}
	svc := monitoring.NewProjectsTimeSeriesService(s)
	call := svc.Query("projects/"+c.projectID, &monitoring.QueryTimeSeriesRequest{
		Query: q,
	})
	return call.Do()
}
