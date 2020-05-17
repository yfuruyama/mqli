package main

import (
	"context"
	"google.golang.org/api/monitoring/v3"
)

type Client struct {
	projectID string
	svc       *monitoring.ProjectsTimeSeriesService
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	s, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		projectID: projectID,
		svc: monitoring.NewProjectsTimeSeriesService(s),
	}, nil
}

func (c *Client) Query(q string) (*Result, error) {
	call := c.svc.Query("projects/"+c.projectID, &monitoring.QueryTimeSeriesRequest{
		Query: q,
	})
	resp, err := call.Do()
	// TODO: error detail
	if err != nil {
		//if err, ok := err.(*googleapi.Error); ok {
		//	log.Printf("details: %#v", err.Details)
		//}
		//log.Printf("resp: %#v", resp)
		//log.Printf("err: %#v", err)
		return nil, err
	}

	result := buildQueryResult(resp.TimeSeriesDescriptor, resp.TimeSeriesData)
	return result, nil
}
