package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

func RegisterCollectors(nlogger *zerolog.Logger) {
	cl := NewPolkadotCollector(nlogger)
	prometheus.MustRegister(prometheus.Collector(cl))
}
