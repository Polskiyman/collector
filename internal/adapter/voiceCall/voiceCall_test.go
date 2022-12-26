package voiceCall

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoiceCall_Parse(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantData []VoiceCallData
		wantErr  error
	}{
		{
			name: "example from task",
			path: "test_voiceCall.data",
			wantData: []VoiceCallData{
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
		},
		{
			name:     "bad file path",
			path:     "qwerty",
			wantData: []VoiceCallData{},
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

func Test_createVoiceCallData(t *testing.T) {
	tests := []struct {
		name    string
		data    []string
		wantRes VoiceCallData
		wantErr error
	}{
		{
			name:    "line not contains 8 fields",
			data:    []string{"AT40673Transparentalls;0.62;581;38;10"},
			wantRes: VoiceCallData{},
			wantErr: errLenFields,
		},
		{
			name:    "line is empty",
			data:    []string{},
			wantRes: VoiceCallData{},
			wantErr: errEmptyLine,
		},
		{
			name: "ok line",
			data: []string{"BG;40;609;E-Voice;0.86;160;36;5"},
			wantRes: VoiceCallData{
				Country:             "BG",
				Bandwidth:           "40",
				ResponseTime:        "609",
				Provider:            "E-Voice",
				ConnectionStability: 0.86,
				TTFB:                160,
				VoicePurity:         36,
				MedianOfCallsTime:   5,
			},
			wantErr: nil,
		},
		{
			name:    "incorrect country code",
			data:    []string{"B1;40;609;E-Voice;0.86;160;36;5"},
			wantRes: VoiceCallData{},
			wantErr: errInvalidCountry,
		},
		{
			name:    "incorrect provider",
			data:    []string{"BL;58;930;E-Voic;0.65;738;83;52"},
			wantRes: VoiceCallData{},
			wantErr: errInvalidProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := createVoiceCallData(tt.data)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, d)
		})
	}
}
