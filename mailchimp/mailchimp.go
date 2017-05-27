package mailchimp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	memberEndpoint string
	apiKey         string
}

type Config struct {
	APIKey string
	ListId string
	URL    string
}

type Member struct {
	EmailAddress string `json:"email_address"`
	Status       string `json:"status"`
}

func New(cfg Config) (client Client, err error) {
	if cfg.APIKey == "" || cfg.ListId == "" || cfg.URL == "" {
		err = errors.Errorf(
			"APIKey, ListID and URL must be all specified.")
		return
	}

	client.apiKey = cfg.APIKey
	client.memberEndpoint = fmt.Sprintf(
		"%s/3.0/lists/%s/members", cfg.URL, cfg.ListId)
	return
}

func (c *Client) Subscribe(email string) (err error) {
	if email == "" {
		err = errors.Errorf(
			"Can't subscribe with empty email")
		return
	}

	jsonBytes, err := json.Marshal(&Member{
		EmailAddress: email,
		Status:       "subscribed",
	})
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't create json for posting to member subscription (email=%s)",
			email)
		return
	}

	req, err := http.NewRequest("POST", c.memberEndpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't create request to perform member subscription")
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.SetBasicAuth("anyuser", c.apiKey)
			req.Header.Add("Content-Type", "application/json")
			return nil
		},
	}

	req.SetBasicAuth("anyuser", c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrapf(err,
			"Request to subscribe email (%s) failed",
			email)
		return
	}

	if resp.StatusCode > 299 {
		err = errors.Errorf(
			"Request wasn't successful (%d - %s)",
			resp.StatusCode, resp.Status)
		return
	}

	return
}
