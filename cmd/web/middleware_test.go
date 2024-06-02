package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mhrdini/snippetbox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	// Recorder that will record the http.Response returned from a handler it is passed into
	rr := httptest.NewRecorder()

	// Dummy request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock handler to pass to secureHeaders
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Fatal(err)
		}
	})

	// Pass mock HTTP handler to secureHeaders, calling ServeHTTP() method to execute with
	// http.ResponseRecorder and dummy *http.Request
	secureHeaders(next).ServeHTTP(rr, r)

	// Get the recorded response
	rs := rr.Result()

	tests := []struct {
		header string
		want   string
	}{
		{
			header: "Content-Security-Policy",
			want:   "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			header: "Referrer-Policy",
			want:   "origin-when-cross-origin",
		},
		{
			header: "X-Content-Type-Options",
			want:   "nosniff",
		},
		{
			header: "X-Frame-Options",
			want:   "deny",
		},
		{
			header: "X-XSS-Protection",
			want:   "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			got := rs.Header.Get(tt.header)
			assert.Equal(t, got, tt.want)
		})
	}

	t.Run("StatusOK", func(t *testing.T) {
		assert.Equal(t, rs.StatusCode, http.StatusOK)
	})

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
