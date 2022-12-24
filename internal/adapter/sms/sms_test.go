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
				{Country: "BL",
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

			assert.Len(t, s.Data, 2)
			assert.Equal(t, s.Data, tt.wantData)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}
