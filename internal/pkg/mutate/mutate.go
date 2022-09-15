package mutate

import (
	"fmt"

	promclient "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

// AppendLabelToMetrics appends a label to the metrics in order to indicate which container the given metric came from.
func AppendLabelToMetrics(labelName string, containerName string, metricFamilies map[string]*promclient.MetricFamily) error {

	if len(labelName) == 0 {
		return fmt.Errorf("empty container label name input")
	}
	if len(containerName) == 0 {
		return fmt.Errorf("empty container name input")
	}
	if len(metricFamilies) == 0 {
		return fmt.Errorf("empty MetricFamily map input")
	}

	for _, mf := range metricFamilies {
		metrics := mf.GetMetric()

		for _, m := range metrics {
			label := promclient.LabelPair{
				Name:  proto.String(labelName),
				Value: proto.String(containerName),
			}
			m.Label = append(m.Label, &label)
		}
	}
	return nil
}
