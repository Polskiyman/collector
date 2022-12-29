package mms

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMms_TestFunc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data" {
			t.Errorf("Expected to request '/data', got: %s", r.URL.Path)
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"test": 2}`))
	}))
	defer server.Close()

	m := New(server.URL)
	v, err := m.TestFunc()
	if v.Test != 1 {
		t.Errorf(`Expected v {Test: 1}, got: %+v`, v)
	}
	if err != nil {
		t.Errorf(`Expected err nil, got: %v`, err)
	}
}
