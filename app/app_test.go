package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestProbeHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest(http.MethodGet, "/liveness", nil)
	assert.Nil(t, err, "should be nil")

	app := &App{}
	app.Init()

	// Create a response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.handleLiveness)

	//Pass our request in
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "status should be ok")

	// Check the response body is what we expect.
	expected := `{}`
	assert.Equal(t, expected, rr.Body.String(), "body didn't match")
}

func TestShortenHandler(t *testing.T) {
	// Create a request to pass to our handler
	bodyJSON, err := json.Marshal(Request{Link: "testlink"})
	assert.Nil(t, err, "should be nil")

	buf := bytes.NewBuffer(bodyJSON)
	req, err := http.NewRequest(http.MethodPost, "/shorten", buf)
	assert.Nil(t, err, "should be nil")

	app := &App{}
	app.Init()

	// Create a response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.handleShorten)

	//Pass our request in
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "status should be ok")

	// Check the response body is what we expect.
	expected := "\"link\":\"localhost:8090"
	assert.Contains(t, rr.Body.String(), expected, "Body didn't contain expected results")
}

func TestRedirectSuccess(t *testing.T) {
	// Create a request to pass to our handler
	bodyJSON, _ := json.Marshal(Request{Link: "testlink"})

	buf := bytes.NewBuffer(bodyJSON)
	req, _ := http.NewRequest(http.MethodPost, "/shorten", buf)

	app := &App{}
	app.Init()

	// Create a response recorder
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/shorten", app.handleShorten)
	r.PathPrefix("/").HandlerFunc(app.handleRedirect)

	//Pass our request in
	r.ServeHTTP(rr, req)
	var shortenedResp Response
	json.Unmarshal(rr.Body.Bytes(), &shortenedResp)

	shortURL := strings.Split(shortenedResp.Link, "localhost:8090")[1]
	req2, err := http.NewRequest("GET", shortURL, nil)
	req2.RequestURI = shortURL
	assert.Nil(t, err, "should be nil")

	//Pass our request in
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req2)
	assert.Equal(t, http.StatusMovedPermanently, rr.Code, "status mismatch")

	// Check the response body is what we expect.
	expected := "Moved Permanently"
	assert.Contains(t, rr.Body.String(), expected, "body didn't match")
}

func TestShortenHandlerErrorsWithInvalidJSON(t *testing.T) {
	// Create a request to pass to our handler
	bodyJSON := "{"
	buf := bytes.NewBufferString(bodyJSON)
	req, err := http.NewRequest(http.MethodPost, "/shorten", buf)
	assert.Nil(t, err, "should be nil")

	app := &App{}
	app.Init()

	// Create a response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.handleShorten)

	//Pass our request in
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "status mismatch")

	// Check the response body is what we expect.
	expected := "{\"error\":\"Error parsing payload\"}\n"
	assert.Equal(t, expected, rr.Body.String(), "body didn't match")
}

func TestRedirectHandlerErrorsWithInvalidLink(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/nolinkhere", nil)
	assert.Nil(t, err, "should be nil")

	app := &App{}
	app.Init()

	// Create a response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.handleRedirect)

	//Pass our request in
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code, "status mismatch")

	// Check the response body is what we expect.
	expected := "{\"error\":\"Invalid link\"}\n"
	assert.Equal(t, expected, rr.Body.String(), "body didn't match")
}
