package glue

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Version = "devel"

	gauge *prometheus.GaugeVec
)

func init() {
	gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "tigris_gtm",
		Subsystem: "glue",
		Name:      "version",
		Help:      "The version of glue in use.",
	}, []string{"version"})

	gauge.WithLabelValues(Version).Inc()
}
