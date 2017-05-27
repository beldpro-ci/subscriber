package mailchimp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/beldpro-ci/subscriber/mailchimp"
	"github.com/stretchr/testify/assert"
)

var ts *httptest.Server

func memberMockHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(r)
}

func TestMain(m *testing.M) {
	ts = httptest.NewServer(http.HandlerFunc(memberMockHandler))
	defer ts.Close()
	os.Exit(m.Run())
}

func TestSubscribeUser_requestIsAuthenticated(t *testing.T) {
	mc, err := mailchimp.New(mailchimp.Config{
		ListId: "123",
		APIKey: "key",
		URL:    ts.URL,
	})
	assert.NoError(t, err)

	err = mc.Subscribe("email")
	assert.NoError(t, err)
}

func TestSubscribeUser_containsData(t *testing.T) {
}

func TestSubscribeUser_containsContentTypeHeader(t *testing.T) {
}
