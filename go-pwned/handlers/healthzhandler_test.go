package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthzHandlerShouldReturnOk(t *testing.T) {
	log := log.New(os.Stdout, "go_test", log.LstdFlags)
	handler := NewHealthzHandler(log)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status: '%d', Got: '%d'\n", http.StatusOK, response.Code)
	}

}
