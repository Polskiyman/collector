package mms

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMms_Fetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "/", r.URL.Path)
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET '/', got: %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		content, _ := os.ReadFile("mms.json")
		// TODO: check error
		_, _ = w.Write(content)
	}))
	defer server.Close()

	mmsAdapter := New(server.URL)
	err := mmsAdapter.Fetch()
	assert.Nil(t, err)
	if err != nil {
		t.Errorf(`Expected err nil, got: %v`, err)
	}

	expectedData := []MmsData{
		{
			Country:      "RU",
			Provider:     "Kildy",
			Bandwidth:    "3",
			ResponseTime: "511",
		},
	}
	assert.Equal(t, expectedData, mmsAdapter.Data)
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
			var m Mms
			m.filterResponse(tt.data)
			assert.Equal(t, tt.wantRes, m.Data)
		})
	}
}
