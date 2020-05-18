package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/monitoring/v3"
	"net/http"
	"strings"
)

const scope = "https://www.googleapis.com/auth/cloud-platform"

type Client struct {
	projectID string
	client    *http.Client
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	// NOTE: We dont' use monitoring.ProjectsTimeSeriesService for API client
	// as it only returns limited error information.
	c, err := google.DefaultClient(ctx, scope)
	if err != nil {
		return nil, err
	}

	return &Client{
		projectID: projectID,
		client:    c,
	}, nil
}

type QueryRequest struct {
	Query string `json:"query"`
}

type ErrorDetail struct {
	Type    string             `json:"@type"`
	Summary string             `json:"errorSummary"`
	Errors  []ErrorDetailError `json:"errors"`
}

type ErrorDetailError struct {
	Message string      `json:"message"`
	Locator interface{} `json:"locator"` // TODO
}

func (c *Client) Query(q string) (*Result, error) {
	b, err := json.Marshal(&QueryRequest{q})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://monitoring.googleapis.com/v3/projects/%s/timeSeries:query", c.projectID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := googleapi.CheckResponse(resp); err != nil {
		if err, ok := err.(*googleapi.Error); ok {
			return nil, convertError(err)
		}
		return nil, err
	}

	// tsResp contains only subset of monitoring.QueryTimeSeriesResponse
	var tsResp monitoring.QueryTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&tsResp); err != nil {
		return nil, err
	}

	return buildQueryResult(tsResp.TimeSeriesDescriptor, tsResp.TimeSeriesData), nil
}

func convertError(original *googleapi.Error) error {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(original.Details); err != nil {
		return original
	}

	var details []ErrorDetail
	if err := json.Unmarshal(b.Bytes(), &details); err != nil {
		return original
	}

	var messages []string
	for _, detail := range details {
		for _, error := range detail.Errors {
			messages = append(messages, error.Message)
		}
	}

	// TODO: Show locator
	return fmt.Errorf("code: %d, message: %q, detail: %q", original.Code, original.Message, strings.Join(messages, ", "))
}
