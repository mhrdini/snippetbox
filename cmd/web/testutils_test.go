package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication() *application {
	return &application{
		// used in the errorLog and recoverPanic middleware used across all routes
		// so we create dummy loggers so those functions won't panic
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add cookie jar to test server client so that any response cookies will be stored and sent with
	// subsequent requests made using this client
	ts.Client().Jar = jar

	// Disable redirect-following for the test server client.
	// It forces client to return received response from a call to 3xx, i.e. redirect requests
	// by always returning a http.ErrUseLastResponse error
	ts.Client().CheckRedirect = func(r *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, path string) (int, http.Header, string) {
	// use the test server Client to send requests
	// the client can be configured to tweak its behaviour
	rs, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
