package main

import (
	"net/http"
	"testing"

	"github.com/mhrdini/snippetbox/internal/assert"
	"github.com/mhrdini/snippetbox/internal/models/mocks"
)

func TestPing(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	type result struct {
		status int
		body   string
	}

	tests := []struct {
		name string
		path string
		want result
	}{
		{
			name: "Valid ID",
			path: "/snippet/view/1",
			want: result{
				status: http.StatusOK,
				body:   mocks.MockSnippet.Content,
			},
		},
		{
			name: "Non-existent ID",
			path: "/snippet/view/2",
			want: result{
				status: http.StatusNotFound,
			},
		},
		{
			name: "Negative ID",
			path: "/snippet/view/-1",
			want: result{
				status: http.StatusNotFound,
			},
		},
		{
			name: "Decimal ID",
			path: "/snippet/view/1.23",
			want: result{
				status: http.StatusNotFound,
			},
		},
		{
			name: "String ID",
			path: "/snippet/view/foo",
			want: result{
				status: http.StatusNotFound,
			},
		},
		{
			name: "Empty ID",
			path: "/snippet/view/",
			want: result{
				status: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, _, body := ts.get(t, tt.path)
			assert.Equal(t, status, tt.want.status)
			if tt.want.body != "" {
				assert.StringContains(t, body, tt.want.body)
			}
		})
	}

}
