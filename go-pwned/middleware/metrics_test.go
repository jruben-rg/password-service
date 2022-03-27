package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jruben-rg/password-service/go-pwned/metric"
)

type TestMetricService struct {
	isInvoked bool
}

func (tm *TestMetricService) SaveMetrics(mi *metric.MetricInfo) {
	tm.isInvoked = true
}

type TestHandler struct {
	isServeHTTPInvoked bool
}

func (th *TestHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	th.isServeHTTPInvoked = true
}

func TestWriteHeaderShouldWriteStatusCodeAndWriteHeader(t *testing.T) {
	tests := []struct {
		scenario         string
		code             int
		expectedRespCode int
	}{
		{
			scenario:         "Should set status code to 200 (Ok)",
			code:             http.StatusOK,
			expectedRespCode: http.StatusOK,
		},
		{
			scenario:         "Should set status code to 400 (BadRequest)",
			code:             http.StatusBadRequest,
			expectedRespCode: http.StatusBadRequest,
		},
		{
			scenario:         "Should set status code to 500 (Internal Server Error)",
			code:             http.StatusInternalServerError,
			expectedRespCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {

		response := httptest.NewRecorder()
		statusRespWriter := NewStatusResponseWriter(response)
		statusRespWriter.WriteHeader(test.code)

		if response.Code != test.expectedRespCode {
			t.Errorf("Scenario '%s' Expected: '%d', Got: '%d'\n", test.scenario, test.expectedRespCode, response.Code)
		}

		if statusRespWriter.code != test.expectedRespCode {
			t.Errorf("Scenario '%s' Expected: '%d', Got: '%d'\n", test.scenario, test.expectedRespCode, statusRespWriter.code)
		}
	}

}

func TestShouldRecordInvokeSaveMetrics(t *testing.T) {

	testMetricSvc := TestMetricService{}
	testHandler := TestHandler{}
	handler := Metrics(&testMetricSvc, &testHandler)

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if testMetricSvc.isInvoked == false {
		t.Errorf("SaveMetrics method wasn't invoked")
	}

	if testHandler.isServeHTTPInvoked == false {
		t.Errorf("ServeHTTP method on handler wasnt invoked")
	}

}
