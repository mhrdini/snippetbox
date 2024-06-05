package main

import (
	"net/http"
	"net/url"
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

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Get CSRF token by making call to GET /user/signup
	_, _, body := ts.get(t, "/user/signup")
	// t.Logf("%v", body)
	csrfToken := extractCSRFToken(t, body)
	// t.Logf("CSRF token is: %q", csrfToken)

	const (
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
		validName     = mocks.ValidName
		validPassword = mocks.ValidPassword
		validEmail    = mocks.ValidEmail
		dupeEmail     = mocks.DupeEmail
	)

	type user struct {
		name     string
		email    string
		password string
	}

	type result struct {
		status  int
		formTag string
	}

	tests := []struct {
		name      string
		user      user
		csrfToken string
		want      result
	}{
		{
			name: "Valid submission",
			user: user{
				name:     validName,
				email:    validEmail,
				password: validPassword,
			},
			csrfToken: csrfToken,
			want: result{
				status: http.StatusSeeOther,
			},
		},
		{
			name: "Short password",
			user: user{
				name:     validName,
				email:    validEmail,
				password: "pa$$",
			},
			csrfToken: csrfToken,
			want: result{
				status:  http.StatusUnprocessableEntity,
				formTag: formTag,
			},
		}, {
			name: "Duplicate email",
			user: user{
				name:     validName,
				email:    dupeEmail,
				password: validPassword,
			},
			csrfToken: csrfToken,
			want: result{
				status:  http.StatusUnprocessableEntity,
				formTag: formTag,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create FormData based on test data
			form := url.Values{}
			form.Add("name", tt.user.name)
			form.Add("email", tt.user.email)
			form.Add("password", tt.user.password)
			form.Add("csrf_token", tt.csrfToken)

			status, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, status, tt.want.status)
			if tt.want.formTag != "" {
				assert.StringContains(t, body, tt.want.formTag)
			}
		})
	}

}
