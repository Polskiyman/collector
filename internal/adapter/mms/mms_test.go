package mms

import (
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
	if len(m.Data) < 1 {
		t.Errorf(`Expected m.Date is nil, got: `)
	}
	if err != nil {
		t.Errorf(`Expected err nil, got: %v`, err)
	}
}
