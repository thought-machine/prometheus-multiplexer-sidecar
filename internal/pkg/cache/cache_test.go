package cache

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_cache "github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/cache/mocks"
)

func TestGetAndInvalidate(t *testing.T) {
	testCases := []struct {
		name          string
		key           string
		ok            bool
		expectedValue []byte
	}{
		{
			"test if returns nil when couldn't find the associated key",
			"random",
			false,
			nil,
		},
		{
			"test if returns correct metric when find the associated key",
			"container1",
			true,
			[]byte("this is metric"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mc := mock_cache.NewMockCache(ctr)
			mc.EXPECT().LoadAndDelete(tc.key).Return(tc.expectedValue, tc.ok)

			cache := RawMetricCache{mc}
			metric, ok := cache.GetAndInvalidate(tc.key)
			assert.Equal(t, tc.expectedValue, metric)
			assert.Equal(t, tc.ok, ok)
		})
	}
}

func TestSet(t *testing.T) {
	testCases := []struct {
		name   string
		key    string
		metric []byte
	}{
		{
			"test if it can handle empty input",
			"random",
			make([]byte, 0),
		},
		{
			"test if it can handle correct input",
			"container1",
			[]byte("this is metric"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mc := mock_cache.NewMockCache(ctr)
			if len(tc.metric) != 0 {
				mc.EXPECT().Store(tc.key, tc.metric)
			}
			cache := RawMetricCache{mc}
			cache.Set(tc.key, tc.metric)
		})
	}
}
