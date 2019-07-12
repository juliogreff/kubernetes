/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package metrics provides abstractions for registering which metrics
// to record.
package metrics

import (
	"net/url"
	"sync"
	"time"
)

var registerMetrics sync.Once

// LatencyMetric observes client latency partitioned by verb and url.
type LatencyMetric interface {
	Observe(verb string, u url.URL, latency time.Duration)
}

// ResultMetric counts response codes partitioned by method and host.
type ResultMetric interface {
	Increment(code string, method string, host string)
}

// ThrottleMetric counts requests that were throttled.
type ThrottleMetric interface {
	Increment(verb string, u url.URL)
}

// ThrottleLatencyMetric observes client latency introduced by throttling.
type ThrottleLatencyMetric interface {
	Observe(verb string, u url.URL, latency time.Duration)
}

type RegisterMetrics struct {
	RequestLatency         LatencyMetric
	RequestResult          ResultMetric
	RequestThrottle        ThrottleMetric
	RequestThrottleLatency ThrottleLatencyMetric
}

var RegisteredMetrics = RegisterMetrics{
	RequestLatency:         noopLatency{},
	RequestResult:          noopResult{},
	RequestThrottle:        noopThrottle{},
	RequestThrottleLatency: noopThrottleLatency{},
}

var (
	// RequestLatency is the latency metric that rest clients will update.
	RequestLatency LatencyMetric = noopLatency{}
	// RequestResult is the result metric that rest clients will update.
	RequestResult ResultMetric = noopResult{}
	// RequestThrottle is the throttling metric that rest clients will update.
	RequestThrottle ThrottleMetric = noopThrottle{}
	// RequestThrottleLatency is the throttling metric metric that rest clients will update.
	RequestThrottleLatency ThrottleLatencyMetric = noopThrottleLatency{}
)

// Register registers metrics for the rest client to use. This can
// only be called once.
func Register(r RegisterMetrics) {
	registerMetrics.Do(func() {
		if r.RequestLatency != nil {
			RequestLatency = r.RequestLatency
		}
		if r.RequestResult != nil {
			RequestResult = r.RequestResult
		}
		if r.RequestThrottle != nil {
			RequestThrottle = r.RequestThrottle
		}
		if r.RequestThrottleLatency != nil {
			RequestThrottleLatency = r.RequestThrottleLatency
		}
	})
}

type noopLatency struct{}

func (noopLatency) Observe(string, url.URL, time.Duration) {}

type noopResult struct{}

func (noopResult) Increment(string, string, string) {}

type noopThrottle struct{}

func (noopThrottle) Increment(string, url.URL) {}

type noopThrottleLatency struct{}

func (noopThrottleLatency) Observe(string, url.URL, time.Duration) {}
