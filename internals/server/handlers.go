package srv

import (
	"net/http"
	"time"

	"github.com/cybrarymin/polkadot_exporter/internals/collector"
)

func (exp *Exporter) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"exporter_version": Version,
		"exporter_status":  "running",
	}

	err := JsonWriter(w, data, http.StatusOK, nil)
	if err != nil {
		exp.serverErrorResponse(w, r, err)
	}

}

func (exp *Exporter) statusHandler(w http.ResponseWriter, r *http.Request) {

	eraNumber, _, _, err := collector.GetCurrentEra(exp.api)
	if err != nil {
		exp.serverErrorResponse(w, r, err)
	}

	rewardPoints, err := collector.GetErasRewardPoints(exp.api, eraNumber)
	if err != nil {
		exp.serverErrorResponse(w, r, err)
	}

	data := map[string]interface{}{
		"timestamp":            time.Now(),
		"current_era":          float64(eraNumber),
		"current_reward_point": rewardPointsParser(rewardPoints),
	}

	err = JsonWriter(w, data, http.StatusOK, nil)
	if err != nil {
		exp.serverErrorResponse(w, r, err)
	}
}
