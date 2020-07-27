package promStats

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartStatistics() {
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "ansem_flag_submitted",
			Help: "Counts submitted flags",
		},
		func() float64 {
			return float64(Stats.GetSubmitted())
		}))

	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "ansem_flag_failed",
			Help: "Counts failed flags",
		},
		func() float64 {
			return float64(Stats.GetFailed())
		}))

	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "ansem_flag_duplicated",
			Help: "Counts duplicated flags",
		},
		func() float64 {
			return float64(Stats.GetDuplicated())
		}))

	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":2112", nil)
}
