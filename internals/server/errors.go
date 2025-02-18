package srv

import (
	"net/http"
)

var (
	error500Message = "internal server error"
	error404Message = "requested resource not found"
	error405Message = "method not allowed for the specified endpoint"
)

type envelope map[string]interface{}

func (exp *Exporter) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	nResponse := envelope{
		"error": message,
	}
	err := JsonWriter(w, nResponse, status, nil)
	if err != nil {
		exp.logger.Error().Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (exp *Exporter) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	exp.logger.Error().Err(err).Send()
	exp.errorResponse(w, r, http.StatusInternalServerError, error500Message)
}

func (exp *Exporter) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	exp.errorResponse(w, r, http.StatusNotFound, error404Message)
}

func (exp *Exporter) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	exp.errorResponse(w, r, http.StatusMethodNotAllowed, error405Message)
}
