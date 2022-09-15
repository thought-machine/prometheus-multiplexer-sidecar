package cache

import "sync"

// Cache is the interface for a cache.
type Cache interface {
	LoadAndDelete(key any) (value any, loaded bool)
	Store(key any, value any)
	Load(key any) (value any, ok bool)
}

// RawMetricCache is the cache to store the raw metrics and their corresponding container names.
type RawMetricCache struct {
	cachedMetrics Cache
}

// NewMetricCache returns a new RawMetricCache pointer.
func NewMetricCache() *RawMetricCache {
	return &RawMetricCache{
		&sync.Map{},
	}
}

// GetAndInvalidate get and invalidate the metric if the metric is stored in the cache.
func (c *RawMetricCache) GetAndInvalidate(containerName string) ([]byte, bool) {
	metric, ok := c.cachedMetrics.LoadAndDelete(containerName)

	if !ok {
		return nil, false
	}
	return metric.([]byte), true
}

// Set sets the value of target metric within the metric cache.
func (c *RawMetricCache) Set(containerName string, metrics []byte) {
	if len(metrics) != 0 {
		c.cachedMetrics.Store(containerName, metrics)
	}
}
