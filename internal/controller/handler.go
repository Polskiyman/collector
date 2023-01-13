package controller

import (
	"encoding/json"
	"net/http"

	"collector/internal/service"
)

func HandleConnection(s service.CollectorInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := readBody(r)
		if err != nil {
			response500(w, err)
			return
		}

		result := s.GetSystemData()
		data, err := json.Marshal(&result)
		if err != nil {
			response500(w, err)
		}
		writeResponse(w, http.StatusOK, data)
	}
}
