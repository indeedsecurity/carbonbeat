package carbonclient

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

func parseNotificationsResBody(body []byte) (Notifications, error) {
	var events Notifications
	err := json.Unmarshal(body, &events)
	if err != nil {
		logp.Debug("api", "parseResBody failed on body: %s", body)
	}
	return events, err
}

func parseAuditEventsResBody(body []byte) (AuditEvents, error) {
	var events AuditEvents
	err := json.Unmarshal(body, &events)
	if err != nil {
		logp.Debug("api", "parseResBody failed on body: %s", body)
	}
	return events, err
}

func authenticatedSIEMGet(c *Client, e string) (*http.Response, error) {
	req, err := http.NewRequest("GET", e, nil)
	if err != nil {
		return &http.Response{}, err
	}
	req.Host = c.cfg.SIEMHost
	req.Header.Set("X-Auth-Token", c.cfg.SIEMKey+"/"+c.cfg.SIEMConnectorID)
	req.Header.Set("User-Agent", "Carbonbeat")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logp.Err(err.Error())
		time.Sleep(333 * time.Minute)
		return resp, err
	}

	return resp, nil
}

func authenticatedAPIGet(c *Client, e string) (*http.Response, error) {
	req, err := http.NewRequest("GET", e, nil)
	if err != nil {
		return &http.Response{}, err
	}
	req.Host = c.cfg.APIHost
	req.Header.Set("X-Auth-Token", c.cfg.APIKey+"/"+c.cfg.APIConnectorID)
	req.Header.Set("User-Agent", "Carbonbeat")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logp.Err(err.Error())
		time.Sleep(3 * time.Minute)
		return resp, err
	}

	return resp, nil
}
