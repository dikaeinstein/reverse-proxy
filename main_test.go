package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestProxy struct{}

func (tp TestProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestReverseProxyHandler(t *testing.T) {
	target, _ := url.Parse("https://localhost:8080")
	proxy := TestProxy{}
	rph := ReverseProxyHandler{target, proxy}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	rph.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected OK got %v", response.Code)
	}
}
