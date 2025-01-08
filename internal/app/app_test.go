package app_test

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alextrufmanov/yandexLMSGo/internal/app"
)

func Test(t *testing.T) {
	testCases := []struct {
		requesMethod       string
		requestBody        string
		expectedAnswerBody string
		expectedStatusCode int
	}{
		{
			requesMethod:       "POST",
			requestBody:        "{\"expression\":\"2+2\"}",
			expectedAnswerBody: "{\"result\":\"4\"}",
			expectedStatusCode: 200,
		},
		{
			requesMethod:       "GET",
			requestBody:        "{\"expression\":\"2+2\"}",
			expectedAnswerBody: "404 not found.",
			expectedStatusCode: 404,
		},
		{
			requesMethod:       "POST",
			requestBody:        "{\"expression\":\"2+\"}",
			expectedAnswerBody: "{\"error\":\"Expression is not valid\"}",
			expectedStatusCode: 422,
		},
		{
			requesMethod:       "POST",
			requestBody:        "{\"expression:\"2+2\"}",
			expectedAnswerBody: "{\"error\":\"Internal server error\"}",
			expectedStatusCode: 500,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.requestBody, func(t *testing.T) {
			// req := httptest.NewRequest("POST", "/", bytes.NewBufferString(testCase.requestBody))
			// req.Header.Set("Content-Type", "application/json")
			// w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.requesMethod, "/api/v1/calculate",
				bytes.NewBufferString(testCase.requestBody))
			w := httptest.NewRecorder()
			app.CalcHandler(w, req)
			answerBody := ""
			statusCode := w.Code
			if statusCode == 200 || statusCode == 404 || statusCode == 422 || statusCode == 500 {
				answerBody = strings.TrimSpace(w.Body.String())
			}
			if answerBody != testCase.expectedAnswerBody || statusCode != testCase.expectedStatusCode {
				t.Errorf("RequestBody: \"%s\". Expected (%d) \"%s\", but found (%d)\"%s\".", testCase.requestBody,
					testCase.expectedStatusCode, testCase.expectedAnswerBody, statusCode, answerBody)
			}
		})
	}
}
