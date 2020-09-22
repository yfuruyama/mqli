package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/grpc/codes"
	"os"
	"testing"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

var (
	skipIntegrateTest bool

	testProjectID string
)

const (
	envTestProjectID = "INTEGRATION_TEST_PROJECT_ID"
)

func TestMain(m *testing.M) {
	initialize()
	os.Exit(m.Run())
}

func initialize() {
	if os.Getenv(envTestProjectID) == "" {
		skipIntegrateTest = true
		return
	}
	testProjectID = os.Getenv(envTestProjectID)
}

func generateRandomMetricType() string {
	return fmt.Sprintf("custom.googleapis.com/mqli/test-metric-%d", time.Now().UnixNano())
}

func insertTimeSeries(t *testing.T, projectID string, labels map[string]string, value int64, isCumulative bool) (string, error) {
	t.Helper()

	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return "", err
	}

	metricType := generateRandomMetricType()
	t.Logf("create time series data in %q", metricType)

	now := time.Now()
	var startTime, endTime time.Time
	if isCumulative {
		startTime = now.Add(-1 * time.Second)
		endTime = now
	} else {
		startTime = now
		endTime = now
	}

	var metricKind metricpb.MetricDescriptor_MetricKind
	if isCumulative {
		metricKind = metricpb.MetricDescriptor_CUMULATIVE
	} else {
		metricKind = metricpb.MetricDescriptor_GAUGE
	}
	req := &monitoringpb.CreateTimeSeriesRequest{
		Name: "projects/" + projectID,
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				MetricKind: metricKind,
				Metric: &metricpb.Metric{
					Type:   metricType,
					Labels: labels,
				},
				Resource: &monitoredres.MonitoredResource{
					Type: "global",
				},
				Points: []*monitoringpb.Point{
					{
						Interval: &monitoringpb.TimeInterval{
							StartTime: &timestamp.Timestamp{Seconds: startTime.Unix()},
							EndTime:   &timestamp.Timestamp{Seconds: endTime.Unix()},
						},
						Value: &monitoringpb.TypedValue{
							Value: &monitoringpb.TypedValue_Int64Value{
								Int64Value: value,
							},
						},
					},
				},
			},
		},
	}

	// CreateTimeSeries call often fails, so we try to retry calling it until it succeeds.
	createTimeSeriesWithRetry := func(ctx context.Context) error {
		retryer := gax.OnCodes([]codes.Code{codes.Internal, codes.Unavailable, codes.Unknown}, gax.Backoff{
			Initial:    time.Second,
			Max:        time.Second * 10,
			Multiplier: 2,
		})
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		for {
			if err := c.CreateTimeSeries(ctx, req); err != nil {
				if delay, shouldRetry := retryer.Retry(err); shouldRetry {
					if err := gax.Sleep(ctx, delay); err != nil {
						return err
					}
					continue
				}
				return err
			}
			return nil
		}
	}
	return metricType, createTimeSeriesWithRetry(ctx)
}

func TestQuery(t *testing.T) {
	if skipIntegrateTest {
		t.Skip("integration test is skipped")
	}

	ctx := context.Background()
	client, err := NewClient(ctx, testProjectID)
	if err != nil {
		t.Fatalf("failed to create mqli client: %v", err)
	}

	for _, tt := range []struct {
		desc         string
		labels       map[string]string
		value        int64
		isCumulative bool
		want         *Result
	}{
		{
			desc:         "gauge data",
			labels:       map[string]string{"foo": "bar"},
			value:        123,
			isCumulative: false,
			want: &Result{
				Header: []string{"project_id", "foo", "time", "value"},
				Rows: []Row{
					Row{
						Columns: []string{testProjectID, "bar", "123"},
					},
				},
			},
		},
		{
			desc:         "cumulative data",
			labels:       map[string]string{"hoge": "fuga"},
			value:        1000,
			isCumulative: true,
			want: &Result{
				Header: []string{"project_id", "hoge", "start_time", "end_time", "value"},
				Rows: []Row{
					Row{
						Columns: []string{testProjectID, "fuga", "1000"},
					},
				},
			},
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			metricType, err := insertTimeSeries(t, testProjectID, tt.labels, tt.value, tt.isCumulative)
			if err != nil {
				t.Fatalf("failed to insert time series data: %v", err)
			}

			// Need to wait for time series data becomes visible
			time.Sleep(time.Second * 10)

			query := fmt.Sprintf("fetch global::%s | for 10m", metricType)
			got, err := client.Query(query)
			if err != nil {
				t.Fatalf("failed to query: %v", err)
			}

			// Ignore time value
			opt := cmpopts.IgnoreSliceElements(func(v string) bool {
				if _, err := time.Parse(time.RFC3339Nano, v); err == nil {
					return true
				}
				return false
			})
			if !cmp.Equal(got, tt.want, opt) {
				t.Errorf("diff = %s", cmp.Diff(got, tt.want, opt))
			}
		})
	}
}
