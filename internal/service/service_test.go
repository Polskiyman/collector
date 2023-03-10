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
	"collector/pkg/config"
)

func Test_collector_GetSystemData(t *testing.T) {
	serverMms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile("../adapter/mms/mms.json")
		assert.Nil(t, err)
		_, _ = w.Write(content)
	}))
	defer serverMms.Close()

	serverIncident := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile("../adapter/incident/incident_data.json")
		assert.Nil(t, err)
		_, _ = w.Write(content)
	}))
	defer serverIncident.Close()

	serverSupport := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile("../adapter/support/support_data.json")
		assert.Nil(t, err)
		_, _ = w.Write(content)
	}))
	defer serverSupport.Close()

	tests := []struct {
		name    string
		want    ResultT
		wantErr error
	}{
		{name: "simple test",
			want: ResultT{
				Status: true,
				Data: &ResultSetT{
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
								Country:      "Russian Federation (the)",
								Provider:     "Kildy",
								Bandwidth:    "3",
								ResponseTime: "511",
							},
						},
						[]mms.MmsData{
							mms.MmsData{
								Country:      "Russian Federation (the)",
								Provider:     "Kildy",
								Bandwidth:    "3",
								ResponseTime: "511",
							},
						},
					},
					VoiceCall: []voiceCall.VoiceCallData{
						{
							Country:             "BG",
							Bandwidth:           "40",
							ResponseTime:        "609",
							Provider:            "E-Voice",
							ConnectionStability: 0.86,
							TTFB:                160,
							VoicePurity:         36,
							MedianOfCallsTime:   5,
						},
						{
							Country:             "DK",
							Bandwidth:           "11",
							ResponseTime:        "743",
							Provider:            "JustPhone",
							ConnectionStability: 0.67,
							TTFB:                82,
							VoicePurity:         74,
							MedianOfCallsTime:   41,
						},
					},
					Email: map[string][][]email.EmailData{
						"RU": [][]email.EmailData{
							[]email.EmailData{
								email.EmailData{
									Country:      "RU",
									Provider:     "Yahoo",
									DeliveryTime: 124,
								},
								email.EmailData{
									Country:      "RU",
									Provider:     "Gmail",
									DeliveryTime: 428,
								},
								email.EmailData{
									Country:      "RU",
									Provider:     "MSN",
									DeliveryTime: 463,
								},
							},
							[]email.EmailData{
								email.EmailData{
									Country:      "RU",
									Provider:     "Gmail",
									DeliveryTime: 428,
								},
								email.EmailData{
									Country:      "RU",
									Provider:     "MSN",
									DeliveryTime: 463,
								},
								email.EmailData{
									Country:      "RU",
									Provider:     "Hotmail",
									DeliveryTime: 592,
								},
							},
						},
						"US": [][]email.EmailData{
							[]email.EmailData{
								email.EmailData{
									Country:      "US",
									Provider:     "Orange",
									DeliveryTime: 45,
								},
								email.EmailData{
									Country:      "US",
									Provider:     "MSN",
									DeliveryTime: 124,
								},
								email.EmailData{
									Country:      "US",
									Provider:     "Yahoo",
									DeliveryTime: 305,
								},
							},
							[]email.EmailData{
								email.EmailData{
									Country:      "US",
									Provider:     "MSN",
									DeliveryTime: 124,
								},
								email.EmailData{
									Country:      "US",
									Provider:     "Yahoo",
									DeliveryTime: 305,
								},
								email.EmailData{
									Country:      "US",
									Provider:     "Hotmail",
									DeliveryTime: 391,
								},
							},
						},
					},
					Billing: billing.BillingData{
						CreateCustomer: true,
						Purchase:       true,
						Payout:         false,
						Recurring:      false,
						FraudControl:   true,
						CheckoutPage:   false,
					},
					Support: []int{
						2,
						36,
					},
					Incidents: []incident.IncidentData{
						{
							Topic:  "Wrong SMS delivery time",
							Status: "active",
						},
						{
							Topic:  "Support overloaded because of EU affect",
							Status: "active",
						},
						{
							Topic:  "Billing isn’t allowed in US",
							Status: "closed",
						},
					},
				},
				Error: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Adapters{
				SmsPath:       "../adapter/sms/test_sms.data",
				MmsUrl:        serverMms.URL,
				VoiceCallPath: "../adapter/voiceCall/test_voice_call.data",
				EmailPath:     "../adapter/email/test_email.data",
				BillingPath:   "../adapter/billing/billing_data_test.txt",
				IncidentUrl:   serverIncident.URL,
				SupportUrl:    serverSupport.URL,
			}
			c := New(cfg)
			got := c.GetSystemData()

			assert.Equal(t, tt.want, got)
		})
	}
}
