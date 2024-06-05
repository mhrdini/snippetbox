package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/mhrdini/snippetbox/internal/models/mocks"
)

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)" />`)

type testServer struct {
	*httptest.Server
}

func extractCSRFToken(t *testing.T, body string) string {
	// matches is an array with the entire matched pattern in the first position,
	// and the values of any captured data in the subsequent positions
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	// we need to unescape the string because:
	// - Go's html/template package automatically escapes all dynamically rendered data; and
	// - the CSRF token is a base64 encoded string that might include the + character
	// which will be escaped to &#43
	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal()
	}

	formDecoder := form.NewDecoder()

	// by not setting session store, SCS will default to in-memory storage (suitable for testing
	// purposes)
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		// used in the errorLog and recoverPanic middleware used across all routes
		// so we create dummy loggers so those functions won't panic
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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

func (ts *testServer) postForm(t *testing.T, path string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+path, form) // the only difference between get and postForm
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
