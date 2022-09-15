package parse

import (
	"bytes"
	"fmt"

	promclient "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// Unmarshal accepts the raw metrics and parse them to MetricFamily objects.
func Unmarshal(metrics *bytes.Buffer) (map[string]*promclient.MetricFamily, error) {
	if metrics == nil {
		return nil, fmt.Errorf("empty raw metrics input")
	}
	parser := expfmt.TextParser{}
	mf, err := parser.TextToMetricFamilies(metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw metrics to MetricFamily: %w", err)
	}
	return mf, nil
}

// Marshal accepts the MetricFamily objects and encodes them into raw metrics.
func Marshal(metricFamilies map[string]*promclient.MetricFamily) (*bytes.Buffer, error) {
	if metricFamilies == nil {
		return nil, fmt.Errorf("empty MetricFamily input")
	}
	rawMetrics := &bytes.Buffer{}
	for _, mf := range metricFamilies {

		buff := &bytes.Buffer{}
		_, err := expfmt.MetricFamilyToText(buff, mf)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MetricFamily to raw metrics: %w", err)
		}

		if _, err := rawMetrics.Write(buff.Bytes()); err != nil {
			return nil, err
		}
	}
	return rawMetrics, nil
}
