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

func Test_createSMSData(t *testing.T) {
	tests := []struct {
		name    string
		data    []string
		wantRes SMSData
		wantErr error
	}{
		{
			name:    "line not contains 4 fields",
			data:    []string{"U5;41910;Topol"},
			wantRes: SMSData{},
			wantErr: errLenFields,
		},
		{
			name:    "line is empty",
			data:    []string{},
			wantRes: SMSData{},
			wantErr: errEmptyLine,
		},
		{
			name: "ok line",
			data: []string{"US;36;1576;Rond"},
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
			data:    []string{"F2;9;484;Topolo"},
			wantRes: SMSData{},
			wantErr: errInvalidCountry,
		},
		{
			name:    "incorrect provider",
			data:    []string{"US;36;1576;Rondsd"},
			wantRes: SMSData{},
			wantErr: errInvalidProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := createSMSData(tt.data)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, d)
		})
	}
}
