package carbonclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

// Options provides API host and credentials
type Options struct {
	ConnectorID string
	APIKey      string
	APIHost     string
	Window      string
}

// Client defines the carbonclient
type Client struct {
	cfg        Options
	httpClient http.Client
}

// New returns a carbonclient
func New(opts Options) (Client, error) {
	c := Client{
		cfg: opts,
		httpClient: http.Client{
			// the default client has no timeout and can hang forever
			Timeout: 5 * time.Second,
		},
	}
	return c, nil
}

// Fetch returns a batch of messages since t, time.Location
func (c *Client) Fetch() (Notifications, error) {
	var events Notifications
	endpoint := fmt.Sprintf(
		"%s/integrationServices/v3/notification",
		c.cfg.APIHost)

	resp, err := authenticatedGet(c, endpoint)
	if err != nil {
		return events, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
			logp.Err("Error closing response: %s", err.Error())
		}
	}()

	// lets deal with all of the documented and undocumented status codes for the
	// api response here
	switch code := resp.StatusCode; code {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return events, err
		}
		logp.Debug("api", "API response body: %s", body)
		events, err = parseResBody(body)
		return events, err
	case 400:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logp.Warn(err.Error())
		}
		return events, fmt.Errorf("BAD REQUEST: %s", body)
	case 401:
		return events, errors.New("not authorized to API")
	case 404:
		return events, errors.New("404: API endpoint not found")
	case 429:
		logp.Warn("recieved error 429 (rate limiting). retrying in a bit")
		time.Sleep(time.Minute)
	case 500:
		logp.Warn("500: Internal server error. retrying in a bit")
		time.Sleep(5 * time.Minute)
	default:
		return events, fmt.Errorf("status code: %v - I don't know what this error "+
			"is, it was not documented", code)
	}

	return events, errors.New("something's wrong")
}
