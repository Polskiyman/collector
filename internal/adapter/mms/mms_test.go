package mms

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMms_Fetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data" {
			t.Errorf("Expected to request '/data', got: %s", r.URL.Path)
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		w.WriteHeader(http.StatusOK)
		content, _ := os.ReadFile("mms.json")
		_, _ = w.Write(content)
	}))
	defer server.Close()

	m := New(server.URL)
	err := m.Fetch()
	if err != nil {
		t.Errorf(`Expected err nil, got: %v`, err)
	}
}
func Test_filterResponse(t *testing.T) {
	tests := []struct {
		name    string
		data    []MmsData
		wantRes []MmsData
	}{
		{
			name: "ok line",
			data: []MmsData{
				{
					Country:      "EN",
					Provider:     "Topolo",
					Bandwidth:    "98",
					ResponseTime: "1920",
				},
				{
					Country:      "RU",
					Provider:     "Kildy",
					Bandwidth:    "3",
					ResponseTime: "511",
				},
			},
			wantRes: []MmsData{
				{
					Country:      "RU",
					Provider:     "Kildy",
					Bandwidth:    "3",
					ResponseTime: "511",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Mms
			data.Data = make([]MmsData, 0)
			data.filterResponse(tt.data)
			assert.Equal(t, tt.wantRes, data.Data)
		})
	}
}
