// Copyright 2021 l1b0k
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aliyun

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var OpenAPILatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "aliyun_openapi_latency",
		Help:    "aliyun openapi latency in ms",
		Buckets: []float64{50, 100, 200, 400, 800, 1600, 3200, 6400, 12800, 13800, 14800, 16800, 20800, 28800, 44800},
	},
	[]string{"api", "error"},
)

// MsSince returns milliseconds since start.
func MsSince(start time.Time) float64 {
	return float64(time.Since(start) / time.Millisecond)
}

// RegisterPrometheus register metrics to prometheus server
func RegisterPrometheus() {
	prometheus.MustRegister(OpenAPILatency)
}
