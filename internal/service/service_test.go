package service

import (
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
								Country:      "BL",
								Bandwidth:    "68",
								ResponseTime: "1594",
								Provider:     "Kildy"},
							sms.SMSData{
								Country:      "US",
								Bandwidth:    "36",
								ResponseTime: "1576",
								Provider:     "Rond"}},
						[]sms.SMSData{
							sms.SMSData{
								Country:      "BL",
								Bandwidth:    "68",
								ResponseTime: "1594",
								Provider:     "Kildy"},
							sms.SMSData{
								Country:      "US",
								Bandwidth:    "36",
								ResponseTime: "1576",
								Provider:     "Rond"}}},
					MMS:       [][]mms.MmsData(nil),
					VoiceCall: []voiceCall.VoiceCallData(nil),
					Email:     map[string][][]email.EmailData(nil),
					Billing:   billing.BillingData{CreateCustomer: false, Purchase: false, Payout: false, Recurring: false, FraudControl: false, CheckoutPage: false},
					Support:   []int(nil),
					Incidents: []incident.IncidentData(nil)},
				Error: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New("../adapter/sms/test_sms.data")
			got := c.GetSystemData()

			assert.Equal(t, tt.want, got)
		})
	}
}
