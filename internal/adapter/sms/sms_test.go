package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSms_Parse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantData []SMSData
		wantErr  error
	}{
		{
			name: "example from task",
			path: "test_sms.data",
			wantData: []SMSData{
				{
					Country:      "US",
					Bandwidth:    "36",
					ResponseTime: "1576",
					Provider:     "Rond",
				},
				{
					Country:      "BL",
					Bandwidth:    "68",
					ResponseTime: "1594",
					Provider:     "Kildy",
				},
			},
		},
		{
			name:     "bad file path",
			path:     "qwerty",
			wantData: []SMSData{},
			wantErr:  errBadPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.path)

			err := s.Parse()

			assert.Equal(t, tt.wantData, s.Data)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_CreateSMSData(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantData []SMSData
		wantErr  error
	}{
		{
			name:     "line not contains 4 fields",
			path:     "test_sms.data",
			wantData: []SMSData{},
			wantErr:  errLenFields,
		},
		{
			name:     "line is empty",
			path:     "test_sms.data",
			wantData: []SMSData{},
			wantErr:  errEmptyLine,
		},
		{
			name:     "incorrect country code",
			path:     "test_sms.data",
			wantData: []SMSData{},
			wantErr:  errCountry,
		},
		{
			name:     "incorrect provider",
			path:     "test_sms.data",
			wantData: []SMSData{},
			wantErr:  errProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.path)

			err := s.Parse()

			assert.Equal(t, tt.wantData, s.Data)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_createSMSData(t *testing.T) {
	data := []string{"US;36;1576;Rond"}
	tests := []struct {
		name    string
		wantRes SMSData
		wantErr error
	}{
		{
			name:    "line not contains 4 fields",
			wantRes: SMSData{},
			wantErr: errLenFields,
		},
		{
			name:    "line is empty",
			wantRes: SMSData{},
			wantErr: errEmptyLine,
		},
		{
			name: "Ok line",
			wantRes: SMSData{

				Country:      "US",
				Bandwidth:    "36",
				ResponseTime: "1576",
				Provider:     "Rond",
			},
			wantErr: nil,
		},
		{
			name:    "incorrect country code",
			wantRes: SMSData{},
			wantErr: errCountry,
		},
		{
			name: "Ok Twoline",
			wantRes: SMSData{

				Country:      "BL",
				Bandwidth:    "68",
				ResponseTime: "1594",
				Provider:     "Kildy",
			},
			wantErr: nil,
		},
		{
			name:    "incorrect provider",
			wantRes: SMSData{},
			wantErr: errProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := createSMSData(data)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, d)
		})
	}
}
