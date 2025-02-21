package collector

import (
	gsrpc "github.com/polkadot-go/api/v4"
	"github.com/polkadot-go/api/v4/types"
)

var RpcBackend string

// EraRewardPoints is a custom struct to match the SCALE-encoded data
type EraRewardPoints struct {
	Total       types.U32
	Individuals map[types.AccountID]types.U32
}

func GetCurrentEra(api *gsrpc.SubstrateAPI) (uint32, string, bool, error) {
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

func GetErasRewardPoints(api *gsrpc.SubstrateAPI, eraIndex uint32) (*EraRewardPoints, error) {
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	eraBytes, err := encodeU32(eraIndex)
	if err != nil {
		return nil, err
	}

	// Create storage key for Staking::ErasRewardPoints(eraIndex)
	key, err := types.CreateStorageKey(meta, "Staking", "ErasRewardPoints", eraBytes, nil)
	if err != nil {
		return nil, err
	}

	// 3) Decode into your custom struct
	var rewardPoints EraRewardPoints // <-- this is the custom struct we defined
	ok, err := api.RPC.State.GetStorageLatest(key, &rewardPoints)
	if err != nil {
		return nil, err
	}

	// If there's nothing stored for that era index, return nil
	if !ok {
		return nil, nil
	}

	return &rewardPoints, nil
}
