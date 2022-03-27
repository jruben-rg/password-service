package pwned

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestExecutePwnedRequestShouldTimeoutIfServerTooSlow(t *testing.T) {

	serverHandler := func(rw http.ResponseWriter, r *http.Request) {

		//Server sleeps for 2 seconds before replying
		time.Sleep(2 * time.Second)
		rw.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(serverHandler))
	defer server.Close()

	client := server.Client()
	request := httptest.NewRequest(http.MethodGet, server.URL, nil)
	//request.RequestURI = ""
	// u, _ := url.Parse(server.URL)
	// request.URL = u
	//Request to pwned service timesout after waitinf for 1 second
	_, err := requestPwnedService(client, request, 1*time.Second)
	if err == nil {
		t.Errorf("Was expecting an error for due to timeout\n")
	}

}

func TestRequestPwnedServiceShouldReturnResponseBody(t *testing.T) {

	responseBody := `01F0E818BACDBF3414DBBDFF52586FEF85B:1
			020D91FC5674B1A103E4B09BFC969C418D0:3
			023B4D8499AB737FF8A5044375B43BF04DA:1
			02DF3C11B37C017CE4B17053F5BB7C20FE4:14`

	serverHandler := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(responseBody))
	}

	server := httptest.NewServer(http.HandlerFunc(serverHandler))
	defer server.Close()

	client := server.Client()
	request := httptest.NewRequest(http.MethodGet, server.URL, nil)
	request.RequestURI = ""
	u, _ := url.Parse(server.URL)
	request.URL = u
	response, err := requestPwnedService(client, request, 1*time.Second)
	if err != nil {
		t.Errorf("Got unexpected error: %s\n", err)
	}
	if result := strings.Compare(response, responseBody); result != 0 {
		t.Errorf("Response returned unexpected body")
	}

}

func TestRequestPwnedServiceRespondsWithErrorIfStatusNotOk(t *testing.T) {

	serverHandler := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(serverHandler))
	defer server.Close()

	client := server.Client()
	request := httptest.NewRequest(http.MethodGet, server.URL, nil)

	_, err := requestPwnedService(client, request, 1*time.Second)
	if err == nil {
		t.Error("Was expecting an error from pwned service", err)
	}

}

func TestIsPasswordInDictionary(t *testing.T) {

	tests := []struct {
		scenario       string
		passwordSuffix string
		responseBody   string
		expectedResult bool
	}{
		{
			scenario:       "Should not find password match if not found in response",
			passwordSuffix: "086C7EF93F3EB77AB9A229A2A70CB2AFB3C",
			responseBody: "0784E50CD59416AFA6E9E22DEBDA9603901:5\n" +
				"078957725007D81F20E2354088A04162EC9:1\n" +
				"0790D55E682FDBFDE7DAF7FCAA14BAE6C71:17\n",
			expectedResult: false,
		},
		{
			scenario:       "Should not find password match if empty response",
			passwordSuffix: "078957725007D81F20E2354088A04162EC9",
			responseBody:   "",
			expectedResult: false,
		},
		{
			scenario:       "Should find password match if found in response",
			passwordSuffix: "078957725007D81F20E2354088A04162EC9",
			responseBody: "0784E50CD59416AFA6E9E22DEBDA9603901:5\n" +
				"078957725007D81F20E2354088A04162EC9:1\n" +
				"0790D55E682FDBFDE7DAF7FCAA14BAE6C71:17\n",
			expectedResult: true,
		},
	}

	for _, test := range tests {
		found := isPasswordInDictionary(test.passwordSuffix, test.responseBody)
		if found != test.expectedResult {
			t.Errorf("Scenario '%s'. Got %t. Expected %t.", test.scenario, found, test.expectedResult)
		}
	}
}

func TestIsSecurePasswordRetursOkIfNotEnabled(t *testing.T) {

	serverHandler := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	// Server should return an errof if invoked.
	server := httptest.NewServer(http.HandlerFunc(serverHandler))
	defer server.Close()
	pwned := Pwned{false, 2 * time.Second, server.URL + "/"}

	isSecure, err := pwned.IsSecurePassword("AnyPassw0rd!")
	if isSecure != true {
		t.Errorf("Expected isSecure: '%t', got: '%t'\n", true, isSecure)
	}

	if err != nil {
		t.Errorf("Wasn't expecting an error, got: '%s'\n", err)
	}

}

func TestIsSecurePasswordForScenarios(t *testing.T) {

	tests := []struct {
		scenario          string
		password          string
		pwnedResponseBody string
		expectedSecure    bool
	}{
		{
			scenario: "Should return false if password is not secure",
			password: "Passw0rd",
			pwnedResponseBody: "52417CB061186A45A9D36B636D126D79A83:5\n" +
				"910077770C8340F63CD2DCA2AC1F120444F:4\n" +
				"53772044DE4A4A934FB68DBC55B504804EA:8\n",
			expectedSecure: false,
		},
		{
			scenario: "Should return true if password is secure",
			password: "Passw0rdSec*",
			pwnedResponseBody: "592EFB557DA4C2F1A267CDD11C6E9A77851:1\n" +
				"595A7EE9F2033FE3004AC0B3EDEB01B56AD:3\n" +
				"596DA3F87BDF9A73DB84BA7108C46160613:1\n",
			expectedSecure: true,
		},
	}

	for _, test := range tests {

		serverHandler := func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(test.pwnedResponseBody))
		}

		server := httptest.NewServer(http.HandlerFunc(serverHandler))
		defer server.Close()

		pwned := Pwned{true, 2 * time.Second, server.URL + "/"}

		isSecure, err := pwned.IsSecurePassword(test.password)
		if err != nil {
			t.Errorf("Wasn't expecting error for scenario '%s'. Got %v\n", test.scenario, err)
		}
		if test.expectedSecure != isSecure {
			t.Errorf("Scenario '%s'. Got %t, expected %t\n", test.scenario, isSecure, test.expectedSecure)
		}
	}

}

func TestIsSecurePasswordShouldReturnErrorIfPwnedServiceReturnsError(t *testing.T) {

	serverHandler := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(serverHandler))
	defer server.Close()

	pwned := Pwned{true, 2 * time.Second, server.URL}

	_, err := pwned.IsSecurePassword("AnyPassw0rd")
	if err == nil {
		t.Error("Was expecting error")
	}

}

func TestEncodeSHA1(t *testing.T) {

	tests := []struct {
		scenario string
		inputStr string
		expected string
	}{
		{
			scenario: "Should encode password to sha1 - Scenario 1",
			inputStr: "ThisIsPassword1",
			expected: "0ED4861207CB4977098BBC60DFFDDF7B3A2CDB76",
		},
		{
			scenario: "Should encode password to sha1 - Scenario 2",
			inputStr: "Passw0rd&Symbols",
			expected: "6BDF06F4A8AEEA5534ECC2651EADD90751DF6724",
		},
	}

	for _, test := range tests {
		got := encodeToSHA1(test.inputStr)
		if result := strings.Compare(test.expected, got); result != 0 {
			t.Errorf("Scenario: '%s'. Expected: %s. Got %s.\n", test.scenario, test.expected, got)
		}
	}

}
