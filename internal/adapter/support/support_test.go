package support

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
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile("support_data.json")
		assert.Nil(t, err)
		_, _ = w.Write(content)
	}))
	defer server.Close()

	supportAdapter := New(server.URL)
	err := supportAdapter.Fetch()
	assert.Nil(t, err)

	expectedData := []SupportData{
		{
			Topic:         "SMS",
			ActiveTickets: 3,
		}, {
			Topic:         "MMS",
			ActiveTickets: 9,
		}, {
			Topic:         "Billing",
			ActiveTickets: 0,
		},
	}
	assert.Equal(t, expectedData, supportAdapter.Data)
}
