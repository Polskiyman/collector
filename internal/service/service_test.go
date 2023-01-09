package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"collector/internal/adapter/billing"
	"collector/internal/adapter/email"
	"collector/internal/adapter/incident"
	"collector/internal/adapter/mms"
	"collector/internal/adapter/sms"
	"collector/internal/adapter/voiceCall"
)

func Test_collector_GetSystemData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile("../adapter/mms/mms.json")
		assert.Nil(t, err)
		_, _ = w.Write(content)
	}))
	defer server.Close()

	tests := []struct {
		name    string
		want    ResultT
		wantErr error
	}{
		{name: "simple test",
			want: ResultT{
				Status: false,
				Data: ResultSetT{
					SMS: [][]sms.SMSData{
						[]sms.SMSData{
							sms.SMSData{
								Country:      "Saint Barthélemy",
								Bandwidth:    "68",
								ResponseTime: "1594",
								Provider:     "Kildy"},
							sms.SMSData{
								Country:      "United States of America (the)",
								Bandwidth:    "36",
								ResponseTime: "1576",
								Provider:     "Rond"}},
						[]sms.SMSData{
							sms.SMSData{
								Country:      "Saint Barthélemy",
								Bandwidth:    "68",
								ResponseTime: "1594",
								Provider:     "Kildy"},
							sms.SMSData{
								Country:      "United States of America (the)",
								Bandwidth:    "36",
								ResponseTime: "1576",
								Provider:     "Rond"}}},
					MMS: [][]mms.MmsData{
						[]mms.MmsData{
							mms.MmsData{
								Country:      "RU",
								Provider:     "Kildy",
								Bandwidth:    "3",
								ResponseTime: "511",
							},
						},
						[]mms.MmsData{
							mms.MmsData{
								Country:      "RU",
								Provider:     "Kildy",
								Bandwidth:    "3",
								ResponseTime: "511",
							},
						},
					},
					VoiceCall: []voiceCall.VoiceCallData(nil),
					Email:     map[string][][]email.EmailData(nil),
					Billing: billing.BillingData{
						CreateCustomer: false,
						Purchase:       false,
						Payout:         false,
						Recurring:      false,
						FraudControl:   false,
						CheckoutPage:   false,
					},
					Support:   []int(nil),
					Incidents: []incident.IncidentData(nil)},
				Error: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New("../adapter/sms/test_sms.data", server.URL)
			got := c.GetSystemData()

			assert.Equal(t, tt.want, got)
		})
	}
}
