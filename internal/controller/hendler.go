package controller

import (
	"collector/internal/service"
	"encoding/json"
	"net/http"
)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	_, err := readBody(r)
	if err != nil {
		response500(w, err)
		return
	}

	c := service.New(
		"C:/Users/User/go/src/skillbox-diploma/sms.data",
		"/mms",
		"C:/Users/User/go/src/skillbox-diploma/voice.data",
		"C:/Users/User/go/src/skillbox-diploma/email.data",
		"C:/Users/User/go/src/skillbox-diploma/billing.data",
		"/accendent",
		"/support")
	got := c.GetSystemData()
	data, err := json.Marshal(&got)
	if err != nil {
		response500(w, err)
	}
	writeResponse(w, http.StatusOK, data)
}
