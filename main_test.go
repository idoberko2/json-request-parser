package jparser

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type teststruct struct {
	Str string
}

func getTestServer() *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc("/fake", func(w http.ResponseWriter, r *http.Request) {
		var s teststruct

		if ok := ParseJSONRequest(w, r, &s); ok {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		}
	})
	return httptest.NewServer(r)
}

func TestParseJSONRequest_ok(t *testing.T) {
	srv := getTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/fake", "application/json", bytes.NewBuffer(
		[]byte(`{
			"str": "test string"
		}`),
	))

	assert.Equal(t, 200, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "OK", string(respBody))
}

func TestParseJSONRequest_badly_formed(t *testing.T) {
	srv := getTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/fake", "application/json", bytes.NewBuffer(
		[]byte(`test string`),
	))

	assert.Equal(t, 400, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"error":"Request body contains badly-formed JSON (at position 2)"}`, string(respBody))
}

func TestParseJSONRequest_badly_formed_eof(t *testing.T) {
	srv := getTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/fake", "application/json", bytes.NewBuffer(
		[]byte(`{
			"str": "test string"
		`),
	))

	assert.Equal(t, 400, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"error":"Request body contains badly-formed JSON"}`, string(respBody))
}

func TestParseJSONRequest_unmarshal_error(t *testing.T) {
	srv := getTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/fake", "application/json", bytes.NewBuffer(
		[]byte(`{
			"str": 2
		}`),
	))

	assert.Equal(t, 400, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"error":"Request body contains an invalid value for the \"Str\" field (at position 13)"}`, string(respBody))
}

func TestParseJSONRequest_unknown_field(t *testing.T) {
	srv := getTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/fake", "application/json", bytes.NewBuffer(
		[]byte(`{
			"str": "test string",
			"other field": "shouldn't be here"
		}`),
	))

	assert.Equal(t, 400, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"error":"Request body contains unknown field \"other field\""}`, string(respBody))
}

func TestParseJSONRequest_empty_json(t *testing.T) {
	srv := getTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/fake", "application/json", bytes.NewBuffer(
		[]byte(""),
	))

	assert.Equal(t, 400, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"error":"Request body must not be empty"}`, string(respBody))
}

func TestParseJSONRequest_multiple_objects(t *testing.T) {
	srv := getTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/fake", "application/json", bytes.NewBuffer(
		[]byte(`{
			"str": "test string"
		}
		{
			"str": "test string"
		}`),
	))

	assert.Equal(t, 400, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"error":"Request body must only contain a single JSON object"}`, string(respBody))
}
