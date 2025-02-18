package collector

import "github.com/prometheus/client_golang/prometheus"

type PolkadotCollector struct {
	CurrentEraDesc       *prometheus.Desc
	ErasRewardPointsDesc *prometheus.Desc
}

func NewPolkadotCollector() *PolkadotCollector {
	return &PolkadotCollector{
		CurrentEraDesc: prometheus.NewDesc(
			"node_current_era_index",
			"This is the latest planned era, depending on how the Session pallet queues the validator set, it might be active or not.",
			[]string{"chain"},
			nil,
		),
		ErasRewardPointsDesc: prometheus.NewDesc(
			"node_eras_reward_points",
			" Rewards for the last [Config::HistoryDepth] eras. If reward hasn't been set or has been removed then 0 reward is returned.",
			[]string{"chain", "validator"},
			nil,
		),
	}

}

func (collector *PolkadotCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.CurrentEraDesc
	ch <- collector.ErasRewardPointsDesc

}

func (collector *PolkadotCollector) Collect(ch chan<- prometheus.Metric) {

	// logic of getting data from rpc endpoint
	ch <- prometheus.MustNewConstMetric(
		collector.CurrentEraDesc,
		prometheus.GaugeValue,
		0,
		"node-01",
	)
}
