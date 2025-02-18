package collector

import "github.com/prometheus/client_golang/prometheus"

func RegisterCollectors() {
	cl := NewPolkadotCollector()
	prometheus.MustRegister(prometheus.Collector(cl))
}
