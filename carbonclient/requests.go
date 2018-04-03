package carbonclient

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

func parseResBody(body []byte) (Notifications, error) {
	var events Notifications
	err := json.Unmarshal(body, &events)
	if err != nil {
		logp.Debug("api", "parseResBody failed on body: %s", body)
	}
	return events, err
}

func authenticatedGet(c *Client, e string) (*http.Response, error) {
	req, err := http.NewRequest("GET", e, nil)
	if err != nil {
		return &http.Response{}, err
	}
	req.Host = c.cfg.APIHost
	req.Header.Set("X-Auth-Token", c.cfg.APIKey+"/"+c.cfg.ConnectorID)
	req.Header.Set("User-Agent", "Carbonbeat")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logp.Err(err.Error())
		time.Sleep(15 * time.Minute)
		return resp, err
	}

	return resp, nil
}
