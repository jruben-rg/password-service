package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type TestHandler struct {
	isNil           bool
	serveHTTPCalled bool
}

func (tf *TestHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	tf.serveHTTPCalled = true
}

func TestPasswordHandlerShouldFailWhenRequestCannotBeParsedFromJson(t *testing.T) {
	log := log.New(os.Stdout, "go_test", log.LstdFlags)
	validatePasswordFunc := func(password string) (bool, error) { return false, fmt.Errorf("An error") }
	handler := NewPasswordHandler(log, validatePasswordFunc, nil)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(`{"something": "cnViZW4K"}`))

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v\n", response.Code)
	}
}

func TestPasswordHandlerShouldFailWhenRequestCannotBeDecoded(t *testing.T) {
	log := log.New(os.Stdout, "go_test", log.LstdFlags)
	validatePasswordFunc := func(password string) (bool, error) { return false, fmt.Errorf("An error") }
	handler := NewPasswordHandler(log, validatePasswordFunc, nil)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(`{"password": "---"}`))

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v\n", response.Code)
	}
}

func TestPasswordHandlerShouldFailWhenPasswordValidationFails(t *testing.T) {
	log := log.New(os.Stdout, "go_test", log.LstdFlags)
	validatePasswordFunc := func(password string) (bool, error) { return false, fmt.Errorf("An error") }
	handler := NewPasswordHandler(log, validatePasswordFunc, nil)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(`{"password": "cnViZW4K"}`))

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v\n", response.Code)
	}
}

func TestPasswordHandlerShouldPassWhenPasswordValidationSucceeds(t *testing.T) {
	log := log.New(os.Stdout, "go_test", log.LstdFlags)
	validatePasswordFunc := func(password string) (bool, error) { return true, nil }
	handler := NewPasswordHandler(log, validatePasswordFunc, nil)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(`{"password": "cnViZW4K"}`))

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("Expected BadRequest got %v\n", response.Code)
	}
}

func TestPasswordHandlerShouldInvokeNextValidatorIfNotNil(t *testing.T) {

	log := log.New(os.Stdout, "gopwned_test", log.LstdFlags)

	tests := []struct {
		scenario           string
		nextHandler        TestHandler
		expectedHTTPCalled bool
	}{
		{
			scenario:           "Should not call handler if handler is nil",
			nextHandler:        TestHandler{isNil: true},
			expectedHTTPCalled: false,
		},
		{
			scenario:           "Should invoke handler if handler is not nil",
			nextHandler:        TestHandler{isNil: false},
			expectedHTTPCalled: true,
		},
	}

	for _, test := range tests {

		validatePasswordFunc := func(password string) (bool, error) { return true, nil }

		var handler http.Handler
		if test.nextHandler.isNil {
			handler = NewPasswordHandler(log, validatePasswordFunc, nil)
		} else {
			handler = NewPasswordHandler(log, validatePasswordFunc, &test.nextHandler)
		}

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(`{"password": "cnViZW4K"}`))
		handler.ServeHTTP(response, request)

		if test.expectedHTTPCalled != test.nextHandler.serveHTTPCalled {
			t.Errorf("scenario '%s'. ServeHTTP method should not have been called. Got: %t, expected %t\n", test.scenario, test.nextHandler.serveHTTPCalled, test.expectedHTTPCalled)
		}

	}

}
