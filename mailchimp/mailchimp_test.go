package mailchimp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beldpro-ci/subscriber/mailchimp"
	"github.com/stretchr/testify/assert"
)

var ts *httptest.Server

func TestSubscribe(t *testing.T) {
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "anyuser", user)
		assert.Equal(t, "key", pass)

		contentTypeHeader := r.Header.Get("Content-Type")
		assert.NotEmpty(t, contentTypeHeader)
		assert.Equal(t, "application/json", contentTypeHeader)

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var member = new(mailchimp.Member)
		err := decoder.Decode(member)
		assert.NoError(t, err)

		assert.Equal(t, "email", member.EmailAddress)
		assert.Equal(t, "subscribed", member.Status)
	}))
	defer ts.Close()

	mc, err := mailchimp.New(mailchimp.Config{
		ListId: "123",
		APIKey: "key",
		URL:    ts.URL,
	})
	assert.NoError(t, err)
	err = mc.Subscribe("email")
	assert.NoError(t, err)
}
