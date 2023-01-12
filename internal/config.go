package internal

import (
	"encoding/json"
	"os"
)

type Config struct {
	PortService  string `json:"port_service"`
	SmsPath      string `json:"sms_path"`
	MmsUrl       string `json:"mms_url"`
	ViceCallPath string `json:"vice_call_path"`
	EmailPath    string `json:"email_path"`
	BillingPath  string `json:"billing_path"`
	IncidentUrl  string `json:"incident_url"`
	SupportUrl   string `json:"support_url"`
}

func ParseFromFile(path string) (cfg Config, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &cfg)
	return
}
