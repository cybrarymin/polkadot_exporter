package collector

import (
	gsrpc "github.com/polkadot-go/api/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

func RegisterCollectors(nlogger *zerolog.Logger) {
	cl := NewPolkadotCollector(nlogger)
	prometheus.MustRegister(prometheus.Collector(cl))
}

type PolkadotCollector struct {
	CurrentEraDesc       *prometheus.Desc
	ErasRewardPointsDesc *prometheus.Desc
	Logger               *zerolog.Logger
	rpcBackend           string
}

func NewPolkadotCollector(nlogger *zerolog.Logger) *PolkadotCollector {
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
		Logger:     nlogger,
		rpcBackend: RpcBackend,
	}

}

func (collector *PolkadotCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.CurrentEraDesc
	ch <- collector.ErasRewardPointsDesc

}

func (collector *PolkadotCollector) Collect(ch chan<- prometheus.Metric) {

	api, err := gsrpc.NewSubstrateAPI(collector.rpcBackend)
	if err != nil {
		collector.Logger.Error().Err(err).Msg("couldn't establish connection to the rpc backend to scrape metrics. backend metrics won't be available")
		return
	}
	defer api.Client.Close()

	// Getting currentEra and populating new metrifc with chain label and currentEra value from staking module
	eraNumber, chain, _, err := GetCurrentEra(api)
	if err != nil {
		collector.Logger.Error().Err(err).Msg("couldn't get the currentEra from backend")
		return
	}

	ch <- prometheus.MustNewConstMetric(
		collector.CurrentEraDesc,
		prometheus.GaugeValue,
		float64(eraNumber),
		chain,
	)

	// Getting rewardPoints
	rewardPoints, err := GetErasRewardPoints(api, eraNumber)
	if err != nil {
		collector.Logger.Error().Err(err).Msg("couldn't get the rewardPoints from backend")
		return
	}
	if rewardPoints == nil {
		// If there's no data, we can choose to do nothing or log it
		collector.Logger.Warn().Msgf("no reward points found for era %d", eraNumber)
		return
	}

	//  Loop over each validator to expose `node_eras_reward_points`
	for validatorAccountID, points := range rewardPoints.Individuals {
		// Convert validator AccountID to something string-like
		ch <- prometheus.MustNewConstMetric(
			collector.ErasRewardPointsDesc,
			prometheus.GaugeValue,
			float64(points),
			chain,
			validatorAccountID.ToHexString(),
		)
	}
}
