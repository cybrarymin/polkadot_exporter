package srv

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (exp *Exporter) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(exp.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(exp.methodNotAllowedResponse)
	// healthcheck endpoitn
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", exp.healthcheckHandler)

	// metrics endpoint
	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())
	// status endpoint
	router.HandlerFunc(http.MethodGet, "/status", exp.statusHandler)

	return router
}
