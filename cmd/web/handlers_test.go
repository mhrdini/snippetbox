package main

import (
	"net/http"
	"testing"

	"github.com/mhrdini/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {

	app := newTestApplication()

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, "OK")
}
