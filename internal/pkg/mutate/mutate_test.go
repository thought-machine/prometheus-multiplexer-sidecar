package mutate

import (
	"errors"
	"reflect"
	"testing"

	promclient "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

const (
	containerName = "container1"
	labelName     = "multiplexed_container"
)

func TestAppendLabelToMetrics(t *testing.T) {
	testCases := []struct {
		name                   string
		labelName              string
		containerName          string
		metricFamilies         map[string]*promclient.MetricFamily
		err                    error
		expectedMetricFamilies map[string]*promclient.MetricFamily
	}{
		{
			"return an error when input empty label name",
			"",
			containerName,
			nil,
			errors.New("empty container label name input"),
			nil,
		},
		{
			"return an error when input empty container name",
			labelName,
			"",
			nil,
			errors.New("empty container name input"),
			nil,
		},
		{
			"return an error when input empty MetricFamily map",
			labelName,
			containerName,
			nil,
			errors.New("empty MetricFamily map input"),
			nil,
		},
		{
			"get the expected result after appending",
			labelName,
			containerName,

			map[string]*promclient.MetricFamily{
				"metric1": &promclient.MetricFamily{
					Name: proto.String("metric1"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(1.234),
							},
							Label: []*promclient.LabelPair{
								&promclient.LabelPair{
									Name:  proto.String("labelname1"),
									Value: proto.String("val1"),
								},
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
							Label: []*promclient.LabelPair{
								&promclient.LabelPair{
									Name:  proto.String("labelname2"),
									Value: proto.String("val2"),
								},
								&promclient.LabelPair{
									Name:  proto.String("labelname3"),
									Value: proto.String("val3"),
								},
								&promclient.LabelPair{
									Name:  proto.String("labelname4"),
									Value: proto.String("val4"),
								},
							},
						},
					},
				},
			},
			nil,
			map[string]*promclient.MetricFamily{
				"metric1": &promclient.MetricFamily{
					Name: proto.String("metric1"),
					Type: promclient.MetricType_UNTYPED.Enum(),
					Metric: []*promclient.Metric{
						&promclient.Metric{
							Untyped: &promclient.Untyped{
								Value: proto.Float64(1.234),
							},
							Label: []*promclient.LabelPair{
								&promclient.LabelPair{
									Name:  proto.String("labelname1"),
									Value: proto.String("val1"),
								},
								&promclient.LabelPair{
									Name:  proto.String(labelName),
									Value: proto.String(containerName),
								},
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
							Label: []*promclient.LabelPair{
								&promclient.LabelPair{
									Name:  proto.String("labelname2"),
									Value: proto.String("val2"),
								},
								&promclient.LabelPair{
									Name:  proto.String("labelname3"),
									Value: proto.String("val3"),
								},
								&promclient.LabelPair{
									Name:  proto.String("labelname4"),
									Value: proto.String("val4"),
								},
								&promclient.LabelPair{
									Name:  proto.String(labelName),
									Value: proto.String(containerName),
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
			err := AppendLabelToMetrics(tc.labelName, tc.containerName, tc.metricFamilies)
			assert.True(t, reflect.DeepEqual(tc.expectedMetricFamilies, tc.metricFamilies))
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
