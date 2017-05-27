package server_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/beldpro-ci/subscriber/mailchimp"
	"github.com/beldpro-ci/subscriber/server"
	"github.com/stretchr/testify/assert"
)

func TestServer_failsIfNoMailChimpClient(t *testing.T) {
	_, err := server.New(server.Config{})
	assert.Error(t, err)
}

func TestSubscribeHandler_requestFailsWithoutEmail(t *testing.T) {
	s, err := server.New(server.Config{
		MailChimp: &mailchimp.Client{},
		Port:      123,
	})
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "http://localhost:123/subscribe", nil)
	w := httptest.NewRecorder()
	s.SubscribeHandler(w, req)

	resp := w.Result()
	assert.Equal(t, 400, resp.StatusCode)
}

func TestSubscribeHandler_requestSucceedsWithEmail(t *testing.T) {
	var requestsReceived = 0
	var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestsReceived += 1
	}))
	defer ts.Close()

	mc, err := mailchimp.New(mailchimp.Config{
		APIKey: "key",
		ListId: "listid",
		URL:    ts.URL,
	})
	assert.NoError(t, err)

	s, err := server.New(server.Config{
		MailChimp: &mc,
		Port:      123,
	})
	assert.NoError(t, err)

	form := url.Values{}
	form.Add("email", "email")

	req := httptest.NewRequest("POST",
		"http://localhost:123/subscribe", strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	s.SubscribeHandler(w, req)

	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, requestsReceived)
}
