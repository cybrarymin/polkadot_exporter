package srv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/cybrarymin/polkadot_exporter/internals/collector"
)

func JsonWriter(w http.ResponseWriter, data interface{}, status int, headers http.Header) error {
	nBuffer := bytes.Buffer{}
	err := json.NewEncoder(&nBuffer).Encode(data)
	if err != nil {
		return err
	}
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(nBuffer.Bytes())
	return nil
}

func ListenAddrParser(addr string) (schema string, host string, err error) {
	nUrl, err := url.Parse(addr)
	if err != nil {
		return "", "", err
	}
	return nUrl.Scheme, nUrl.Host, nil
}

func rewardPointsParser(rewardPoints *collector.EraRewardPoints) map[string]uint32 {
	var pointsMap map[string]uint32
	if rewardPoints != nil {
		// Transform map[types.AccountID32]types.U32 into map[string]uint32
		pointsMap = make(map[string]uint32, len(rewardPoints.Individuals))
		for accountID, points := range rewardPoints.Individuals {
			// Convert the 32-byte accountID to a hex string
			accountStr := accountID.ToHexString()
			// Convert the types.U32 to a plain Go uint32
			pointsMap[accountStr] = uint32(points)
		}
	} else {
		// No reward points for this era
		return nil
	}
	return pointsMap
}
