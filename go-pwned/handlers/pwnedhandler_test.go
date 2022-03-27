package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type TestPwnedValidator struct {
	ReturnIsSecure bool
	ReturnError    error
	HasBeenInvoked bool
}

func (pv *TestPwnedValidator) TestValidatePassword(password string) (bool, error) {
	pv.HasBeenInvoked = true
	return pv.ReturnIsSecure, pv.ReturnError
}

func TestRetrieveUserPasswordShouldReturnError(t *testing.T) {

	request := httptest.NewRequest("POST", "/validate", nil)

	tests := []struct {
		scenario string
		key      interface{}
		value    interface{}
	}{
		{
			scenario: "Password is not set in the request context",
			key:      nil,
			value:    nil,
		},
		{
			scenario: "Password is empty in request context",
			key:      PwnedContextKey("UserPassword"),
			value:    "",
		},
	}

	for _, test := range tests {

		if test.key != nil {
			context := context.WithValue(request.Context(), test.key, test.value)
			request = request.WithContext(context)
		}

		_, err := getPassword(request)

		if err == nil {
			t.Errorf("Was expecting error for scenario '%s'\n", test.scenario)
		}
	}

}

func TestRetrieveUserPasswordShouldReturnPassword(t *testing.T) {

	password := "Passw0rd"
	request := httptest.NewRequest("POST", "/validate", nil)

	context := context.WithValue(request.Context(), PwnedContextKey("UserPassword"), password)
	request = request.WithContext(context)

	value, err := getPassword(request)

	if err != nil {
		t.Errorf("Wasn't expecting error, Got: '%s'\n", err)
	}

	if result := strings.Compare(password, value); result != 0 {
		t.Errorf("Was expecting password: '%s', Got: '%s'\n", password, value)
	}

}

func TestPwnedHandler(t *testing.T) {

	log := log.New(os.Stdout, "gopwned_test", log.LstdFlags)

	tests := []struct {
		scenario                 string
		password                 interface{}
		expectedResponseCode     int
		validatorIsSecure        bool
		validatorError           error
		expectedValidatorInvoked bool
	}{
		{
			scenario:                 "Should respond with InternalServerError (500) if password is nil",
			password:                 nil,
			expectedResponseCode:     http.StatusInternalServerError,
			validatorIsSecure:        false,
			validatorError:           nil,
			expectedValidatorInvoked: false,
		},
		{
			scenario:                 "Should respond with InternalServerError (500) if password is not set",
			password:                 "",
			expectedResponseCode:     http.StatusInternalServerError,
			validatorIsSecure:        false,
			validatorError:           nil,
			expectedValidatorInvoked: false,
		},
		{
			scenario:                 "Should respond with InternalServerError (500) if error in the pwned validator",
			password:                 "Passw0rd",
			expectedResponseCode:     http.StatusInternalServerError,
			validatorIsSecure:        false,
			validatorError:           fmt.Errorf("test validator error"),
			expectedValidatorInvoked: true,
		},
		{
			scenario:                 "Should respond with BadRequest (400) if password is insecure",
			password:                 "Passw0rd",
			expectedResponseCode:     http.StatusBadRequest,
			validatorIsSecure:        false,
			validatorError:           nil,
			expectedValidatorInvoked: true,
		},
		{
			scenario:                 "Should respond with Ok (200) if password is secure",
			password:                 "Passw0rd",
			expectedResponseCode:     http.StatusOK,
			validatorIsSecure:        true,
			validatorError:           nil,
			expectedValidatorInvoked: true,
		},
	}

	for _, test := range tests {

		validator := TestPwnedValidator{ReturnIsSecure: test.validatorIsSecure, ReturnError: test.validatorError}
		handler := NewPwnedHandler(log, validator.TestValidatePassword)

		//Create new request and response
		request := httptest.NewRequest("POST", "/validate", nil)
		response := httptest.NewRecorder()

		//Set password in request context
		if test.password != nil {
			c := context.WithValue(request.Context(), PwnedContextKey("UserPassword"), test.password)
			request = request.WithContext(c)
		}

		handler.ServeHTTP(response, request)

		// Verify expected response
		if response.Code != test.expectedResponseCode {
			t.Errorf("Scenario '%s'. Expected Response Code: %d. Got: %d.\n", test.scenario, test.expectedResponseCode, response.Code)
		}

		// Verify whether if validator should have been invoked
		if test.expectedValidatorInvoked != validator.HasBeenInvoked {
			t.Errorf("Scenario '%s'. Expected validator invoked: %t, Got: %t\n", test.scenario, test.expectedValidatorInvoked, validator.HasBeenInvoked)
		}
	}

}
