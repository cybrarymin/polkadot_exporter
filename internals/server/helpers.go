package srv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
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
