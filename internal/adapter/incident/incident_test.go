package incident

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMms_Fetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		content, err := os.ReadFile("incident_data.json")
		assert.Nil(t, err)
		_, _ = w.Write(content)
	}))
	defer server.Close()

	incidentAdapter := New(server.URL)
	err := incidentAdapter.Fetch()
	assert.Nil(t, err)

	expectedData := []IncidentData{
		{
			Topic:  "Billing isn’t allowed in US",
			Status: "closed",
		}, {
			Topic:  "Wrong SMS delivery time",
			Status: "active",
		}, {
			Topic:  "Support overloaded because of EU affect",
			Status: "active",
		},
	}
	assert.Equal(t, expectedData, incidentAdapter.Data)
}
