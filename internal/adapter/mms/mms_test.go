package mms

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMms_Fetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "/", r.URL.Path)

		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile("mms.json")
		assert.Nil(t, err)

		_, _ = w.Write(content)
	}))
	defer server.Close()

	mmsAdapter := New(server.URL)
	err := mmsAdapter.Fetch()
	assert.Nil(t, err)

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

func TestMms_GetContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "/", r.URL.Path)

		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusInternalServerError)

		_, _ = w.Write(nil)
	}))
	defer server.Close()

	mmsAdapter := New(server.URL)
	_, err := mmsAdapter.GetContent()
	errWant := fmt.Errorf("not success")
	assert.Equal(t, err, errWant)

}
