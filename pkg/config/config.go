package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	AppUrl   string   `json:"app_url"`
	Adapters Adapters `json:"adapters"`
}

type Adapters struct {
	SmsPath       string `json:"sms_path"`
	MmsUrl        string `json:"mms_url"`
	VoiceCallPath string `json:"voice_call_path"`
	EmailPath     string `json:"email_path"`
	BillingPath   string `json:"billing_path"`
	IncidentUrl   string `json:"incident_url"`
	SupportUrl    string `json:"support_url"`
}

func ParseFromFile(path string) (cfg Config, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &cfg)
	return
}
