package proxy

import (
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	P                    *Proxy
	preProcessHistogram  prometheus.Histogram
	postProcessHistogram prometheus.Histogram
)

func init() {
	rand.Seed(time.Now().UnixNano())

	preProcessHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "pre_process_time_us",
		Help:    "Includes retry times",
		Buckets: []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
	})

	postProcessHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "post_process_time_us",
		Help:    "",
		Buckets: []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
	})
	P = LoadConfig()
}
