package checkbodyplugin_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	checkbody "github.com/soulbalz/checkbodyplugin"
)

var required = true
var not_required = false
var contains = true
var urlDecode = true

func TestBodyMatch(t *testing.T) {
	var requestBody = []byte(`{
		"test1": "testvalue1",
		"test2": "testvalue2",
		"test3": "testvalue3",
		"test4": "value4"
	}`)

	executeTest(t, requestBody, http.StatusOK)
}

func TestBodyOneMatch(t *testing.T) {
	var requestBody = []byte(`{
		"test1": "testvalue1",
		"test2": "testvalue2",
		"test3": "testvalue3",
		"test4": "value4"
	}`)

	executeTest(t, requestBody, http.StatusOK)
}

func TestBodyNotMatch(t *testing.T) {
	var requestBody = []byte(`{
		"test1": "wrongvalue1",
		"test2": "wrongvalue2",
		"test3": "wrongvalue3",
		"test4": "correctvalue4"
	}`)

	executeTest(t, requestBody, http.StatusForbidden)
}

func TestBodyNotRequired(t *testing.T) {
	var requestBody = []byte(`{
		"test1": "testvalue1",
		"test2": "testvalue2",
		"test4": "ue4"
	}`)

	executeTest(t, requestBody, http.StatusOK)
}

func executeTest(t *testing.T, requestBody []byte, expectedResultCode int) {
	cfg := checkbody.CreateConfig()
	cfg.Body = []checkbody.SingleBody{
		{
			Name:      "test1",
			MatchType: string(checkbody.MatchOne),
			Values:    []string{"testvalue1"},
		},
		{
			Name:      "test2",
			MatchType: string(checkbody.MatchOne),
			Values:    []string{"testvalue2"},
			Required:  &required,
		},
		{
			Name:      "test3",
			MatchType: string(checkbody.MatchOne),
			Values:    []string{"testvalue3"},
			Required:  &not_required,
		},
		{
			Name:      "test4",
			MatchType: string(checkbody.MatchOne),
			Values:    []string{"ue4"},
			Required:  &required,
			Contains:  &contains,
		},
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := checkbody.New(ctx, next, cfg, "check-body-plugin")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != expectedResultCode {
		t.Errorf("Unexpected response status code: %d, expected: %d for incoming request Body: %s", recorder.Result().StatusCode, expectedResultCode, requestBody)
	}
}
