package app

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/indeedsecurity/carbonbeat/carbonclient"
	"github.com/indeedsecurity/carbonbeat/config"
)

// Carbonbeat is the parent that provides fields for the methods
type Carbonbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
	cb     carbonclient.Client
}

// New creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	// authenticated carbonblack API client
	cb, err := carbonclient.New(carbonclient.Options{
		ConnectorID: config.API.ID,
		APIKey:      config.API.Key,
		APIHost:     config.API.Host,
	})
	if err != nil {
		return nil, err
	}

	bt := &Carbonbeat{
		done:   make(chan struct{}),
		config: config,
		cb:     cb,
	}
	return bt, nil
}

// Run starts the beater daemon
func (bt *Carbonbeat) Run(b *beat.Beat) error {
	logp.Info("Carbonbeat is running! Hit CTRL-C to stop it.")

	const maxRetryLimit = 3
	failsSinceLastTry := 0
	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			logp.Warn("recieved done signal, shutting down")
			return nil
		case <-ticker.C:
			//Remove all ticks from the ticker channel, do not "catch up"
			for len(ticker.C) > 0 {
				<-ticker.C
			}
		}

		notifications, err := bt.cb.FetchNotifications()
		if err != nil {
			if failsSinceLastTry > maxRetryLimit {
				return err
			}
			failsSinceLastTry++
			logp.Critical("fetching notifications from the API failed, got: ", err)
		} else {
			failsSinceLastTry = 0
		}

		processedNotifications, err := bt.processNotifications(notifications)
		if err != nil {
			logp.Critical("processing notifications failed because of: ", err)
			return nil
		}

		// goes to output
		bt.client.PublishEvents(processedNotifications, publisher.Guaranteed)
		logp.Info("Sent %d events", len(processedNotifications))
		logp.Debug("api", "events sent: %v", processedNotifications)
	}
}

// Stop gets called when libbeat gets a SIGTERM. It sends a message in a channel to
// stop the worker.
func (bt *Carbonbeat) Stop() {
	err := bt.client.Close()
	if err != nil {
		logp.Critical("stopping the beat client failed because of: ", err)
	}
	close(bt.done)
}
