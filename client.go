package main

import (
	"context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/monitoring/v3"
	"log"
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
	// TODO: error detail
	if err != nil {
		if err, ok := err.(*googleapi.Error); ok {
			log.Printf("details: %#v", err.Details)
		}
		log.Printf("resp: %#v", resp)
		log.Printf("err: %#v", err)
		return nil, err
	}

	result := buildQueryResult(resp.TimeSeriesDescriptor, resp.TimeSeriesData)
	return result, nil
}
