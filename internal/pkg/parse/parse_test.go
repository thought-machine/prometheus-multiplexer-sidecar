package parse

import (
	"bytes"
	"errors"
	"math"
	"strings"
	"testing"

	promclient "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestRawMetricsToMetricFamilies(t *testing.T) {
	testCases := []struct {
		name           string
		rawMetric      *bytes.Buffer
		err            error
		metricFamilies map[string]*promclient.MetricFamily
	}{
		{
			"returns an error when the input is empty",
			nil,
			errors.New("empty raw metrics input"),
			nil,
		},
		{
			"return correct MetricFamily map when inputs multiple types of metrics",
			bytes.NewBuffer([]byte(`minimal_metric 1.234
	another_metric -3e3 103948
	no_labels{} 3
	`)),
			nil,
			map[string]*promclient.MetricFamily{
				"minimal_metric": &promclient.MetricFamily{
					Name: proto.String("minimal_metric"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(1.234),
							},
						},
					},
				},
				"another_metric": &promclient.MetricFamily{
					Name: proto.String("another_metric"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(-3e3),
							},
							TimestampMs: proto.Int64(103948),
						},
					},
				},
				"no_labels": &promclient.MetricFamily{Name: proto.String("no_labels"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(3),
							},
						},
					}},
			},
		},

		{
			"return correct MetricFamily map when inputs one type of metrics",
			bytes.NewBuffer([]byte(`# HELP just_histogram A simple test case for histogram metrics.
# TYPE just_histogram histogram
just_histogram_bucket{le="1"} 209
just_histogram_bucket{le="2"} 210
just_histogram_bucket{le="+Inf"} 2693
just_histogram_sum 314159.26535
just_histogram_count 2333
`)),
			nil,
			map[string]*promclient.MetricFamily{
				"just_histogram": &promclient.MetricFamily{
					Name: proto.String("just_histogram"),
					Help: proto.String("A simple test case for histogram metrics."),
					Type: promclient.MetricType_HISTOGRAM.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Histogram: &promclient.Histogram{
								SampleCount: proto.Uint64(2333),
								SampleSum:   proto.Float64(314159.26535),
								Bucket: []*promclient.Bucket{
									&promclient.Bucket{
										UpperBound:      proto.Float64(1),
										CumulativeCount: proto.Uint64(209),
									},
									&promclient.Bucket{
										UpperBound:      proto.Float64(2),
										CumulativeCount: proto.Uint64(210),
									},
									&promclient.Bucket{
										UpperBound:      proto.Float64(math.Inf(+1)),
										CumulativeCount: proto.Uint64(2693),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"returns an error when inputs invalid raw metrics content",
			bytes.NewBuffer([]byte(`
		Random! Just random words that cannot be parsed into MetricFamily.`)),
			errors.New("failed to parse raw metrics to MetricFamily"),
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mf, err := Unmarshal(tc.rawMetric)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, mf, tc.metricFamilies)
		})
	}
}

func TestMetricFamilyToRawMetrics(t *testing.T) {

	testCases := []struct {
		name           string
		rawMetric      *bytes.Buffer
		err            error
		metricFamilies map[string]*promclient.MetricFamily
	}{
		{
			"returns an error when the MetricFamily input is empty",
			nil,
			errors.New("empty MetricFamily input"),
			nil,
		},
		{
			"return correct raw metrics when inputs multiple MetricFamily objects",
			bytes.NewBuffer([]byte(`# TYPE minimal_metric untyped
minimal_metric 1.234
# TYPE another_metric untyped
another_metric -3000 103948
# TYPE no_labels untyped
no_labels 3
`)),
			nil,
			map[string]*promclient.MetricFamily{
				"minimal_metric": &promclient.MetricFamily{
					Name: proto.String("minimal_metric"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(1.234),
							},
						},
					},
				},
				"another_metric": &promclient.MetricFamily{
					Name: proto.String("another_metric"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(-3e3),
							},
							TimestampMs: proto.Int64(103948),
						},
					},
				},
				"no_labels": &promclient.MetricFamily{Name: proto.String("no_labels"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(3),
							},
						},
					}},
			},
		},

		{
			"return correct raw metrics when inputs one MetricFamily object",
			bytes.NewBuffer([]byte(`# HELP just_histogram A simple test case for histogram metrics.
# TYPE just_histogram histogram
just_histogram_bucket{le="1"} 209
just_histogram_bucket{le="2"} 210
just_histogram_bucket{le="+Inf"} 2693
just_histogram_sum 314159.26535
just_histogram_count 2333
`)),
			nil,
			map[string]*promclient.MetricFamily{
				"just_histogram": &promclient.MetricFamily{
					Name: proto.String("just_histogram"),
					Help: proto.String("A simple test case for histogram metrics."),
					Type: promclient.MetricType_HISTOGRAM.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Histogram: &promclient.Histogram{
								SampleCount: proto.Uint64(2333),
								SampleSum:   proto.Float64(314159.26535),
								Bucket: []*promclient.Bucket{
									&promclient.Bucket{
										UpperBound:      proto.Float64(1),
										CumulativeCount: proto.Uint64(209),
									},
									&promclient.Bucket{
										UpperBound:      proto.Float64(2),
										CumulativeCount: proto.Uint64(210),
									},
									&promclient.Bucket{
										UpperBound:      proto.Float64(math.Inf(+1)),
										CumulativeCount: proto.Uint64(2693),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rm, err := Marshal(tc.metricFamilies)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			expectedMetricList := strings.SplitAfter(tc.rawMetric.String(), "\n")
			actualMetricList := strings.SplitAfter(rm.String(), "\n")
			assert.ElementsMatch(t, expectedMetricList, actualMetricList)
		})
	}
}
