package srv

import (
	"net/http"
	"time"
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
	data := map[string]interface{}{
		"timestamp":            time.Now(),
		"current_era":          "MYTESTERA",
		"current_reward_point": "MYTESTREWARDPOINT",
	}
	err := JsonWriter(w, data, http.StatusOK, nil)
	if err != nil {
		exp.serverErrorResponse(w, r, err)
	}
}
