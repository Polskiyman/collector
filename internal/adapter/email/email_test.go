package email

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"collector/pkg/country"
	"collector/pkg/provider"
)

func TestEmail_Parse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantData []EmailData
		wantErr  error
	}{
		{
			name: "example",
			path: "test_email.data",
			wantData: []EmailData{
				{Country: "RU", Provider: "Gmail", DeliveryTime: 428},
				{Country: "RU", Provider: "Yahoo", DeliveryTime: 124},
				{Country: "RU", Provider: "Hotmail", DeliveryTime: 592},
				{Country: "RU", Provider: "MSN", DeliveryTime: 463},
				{Country: "US", Provider: "Yahoo", DeliveryTime: 305},
				{Country: "US", Provider: "Hotmail", DeliveryTime: 391},
				{Country: "US", Provider: "MSN", DeliveryTime: 124},
				{Country: "US", Provider: "Orange", DeliveryTime: 45},
			},
		},
		{
			name:     "bad file path",
			path:     "qwerty",
			wantData: []EmailData{},
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

func Test_createEmailData(t *testing.T) {
	tests := []struct {
		name    string
		data    []string
		wantRes EmailData
		wantErr error
	}{
		{
			name:    "line not contains 3 fields",
			data:    []string{"AT;Yahoo274"},
			wantRes: EmailData{},
			wantErr: errLenFields,
		},
		{
			name:    "line is empty",
			data:    []string{},
			wantRes: EmailData{},
			wantErr: errEmptyLine,
		},
		{
			name: "ok line",
			data: []string{"AT;Hotmail;487"},
			wantRes: EmailData{

				Country:      "AT",
				Provider:     "Hotmail",
				DeliveryTime: 487,
			},
			wantErr: nil,
		},
		{
			name:    "incorrect country code",
			data:    []string{"T;Gmail;511"},
			wantRes: EmailData{},
			wantErr: country.ErrInvalidCountry,
		},
		{
			name:    "incorrect provider",
			data:    []string{"AT;Hotmaill;487"},
			wantRes: EmailData{},
			wantErr: provider.ErrInvalidProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := createEmailData(tt.data)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, d)
		})
	}
}
