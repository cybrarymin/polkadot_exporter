package collector

import (
	gsrpc "github.com/polkadot-go/api/v4"
	"github.com/polkadot-go/api/v4/types"
	"github.com/prometheus/client_golang/prometheus"
)

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
	wsurl := "ws://localhost:9944"
	api, _ := gsrpc.NewSubstrateAPI(wsurl)
	defer api.Client.Close()
	number, chain, _, _ := getCurrentEra(api)

	// logic of getting data from rpc endpoint
	ch <- prometheus.MustNewConstMetric(
		collector.CurrentEraDesc,
		prometheus.GaugeValue,
		float64(number),
		chain,
	)
}

func getCurrentEra(api *gsrpc.SubstrateAPI) (uint32, string, bool, error) {
	chain, err := api.RPC.System.Chain()
	if err != nil {
		return 0, "", false, err
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return 0, "", false, err
	}

	// Create storage key for Staking::CurrentEra
	key, err := types.CreateStorageKey(meta, "Staking", "CurrentEra", nil, nil)
	if err != nil {
		return 0, "", false, err
	}

	// Because currentEra is `Option<u32>`, we decode it into types.OptionU32
	var currentEraOpt types.OptionU32
	ok, err := api.RPC.State.GetStorageLatest(key, &currentEraOpt)
	if err != nil {
		return 0, "", false, err
	}

	// If storage not found or the Option is None, we return false
	if !ok || currentEraOpt.IsNone() {
		return 0, "", false, nil
	}

	// Unwrap the Option<u32> to get the actual value
	_, eraVal := currentEraOpt.Unwrap()
	return uint32(eraVal), string(chain), true, nil
}

// func getErasRewardPoints(api *gsrpc.SubstrateAPI, eraIndex uint32) (*types.Rewa, error) {
// 	meta, err := api.RPC.State.GetMetadataLatest()
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting metadata: %w", err)
// 	}

// 	// Create storage key for Staking::ErasRewardPoints(eraIndex)
// 	// This is a map-like storage: "ErasRewardPoints" => (EraIndex) => EraRewardPoints
// 	// So we encode eraIndex as the second parameter.
// 	eraBytes, err := types.Encode(eraIndex)
// 	if err != nil {
// 		return nil, fmt.Errorf("error encoding eraIndex: %w", err)
// 	}

// 	key, err := types.CreateStorageKey(meta, "Staking", "ErasRewardPoints", eraBytes, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating storage key for erasRewardPoints: %w", err)
// 	}

// 	// Decode into EraRewardPoints
// 	var rewardPoints types.EraRewardPoints
// 	ok, err := api.RPC.State.GetStorageLatest(key, &rewardPoints)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching ErasRewardPoints for era %d: %w", eraIndex, err)
// 	}

// 	if !ok {
// 		// Means there's no reward points stored for that era
// 		return nil, nil
// 	}

// 	return &rewardPoints, nil
// }
