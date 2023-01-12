package controller

import (
	"io"
	"net/http"
)

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func response500(w http.ResponseWriter, err error) {
	writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
}

func writeResponse(w http.ResponseWriter, statusCode int, resp []byte) {
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}
