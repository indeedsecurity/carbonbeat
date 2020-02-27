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
		APIConnectorID:  config.API.ID,
		APIKey:          config.API.Key,
		APIHost:         config.API.Host,
		SIEMConnectorID: config.SIEM.ID,
		SIEMKey:         config.SIEM.Key,
		SIEMHost:        config.SIEM.Host,
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
	siemFailsSinceLastSuccess := 0
	apiFailsSinceLastSuccess := 0
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
		//These could be moved to goroutines, but that isn't currently necessary
		siemErr := bt.FetchAndSendSIEMEvents()
		if siemErr != nil {
			if siemFailsSinceLastSuccess > maxRetryLimit {
				return siemErr
			}
			siemFailsSinceLastSuccess++
			logp.Critical("Fetching SIEM events failed, got: %s", siemErr)
		} else {
			siemFailsSinceLastSuccess = 0
		}

		apiErr := bt.FetchAndSendAPIEvents()
		if apiErr != nil {
			if apiFailsSinceLastSuccess > maxRetryLimit {
				return apiErr
			}
			apiFailsSinceLastSuccess++
			logp.Critical("Fetching API events failed, got: %s", apiErr)
		} else {
			apiFailsSinceLastSuccess = 0
		}
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

//FetchAndSendSIEMEvents fetches the carbonblack events from the SIEM key, marshalls them into a map and
//sends them to the configured output
func (bt *Carbonbeat) FetchAndSendSIEMEvents() error {
	//That endpoint is for carbonblack event notifications
	notifications, err := bt.cb.FetchSIEMEvents("/integrationServices/v3/notification")
	if err != nil {
		logp.Critical("fetching notifications from the API failed, got: %s", err)
		return err
	}
	processedNotifications, err := bt.processNotifications(notifications)
	if err != nil {
		logp.Critical("processing notifications failed because of: %s", err)
		return err
	}
	bt.client.PublishEvents(processedNotifications, publisher.Guaranteed)
	logp.Debug("api", "notification events sent: %v", processedNotifications)
	return nil
}

//FetchAndSendAPIEvents fetches the carbonblack audit log events, marshalls them into a map
//and sends them to the configured output
func (bt *Carbonbeat) FetchAndSendAPIEvents() error {
	//That endpoint is for carbonblack audit log events
	events, err := bt.cb.FetchAPIEvents("/integrationServices/v3/auditlogs")
	if err != nil {
		logp.Critical("fetching audit events from the API failed, got: %s", err)
		return err
	}
	processedEvents, err := bt.processAuditEvents(events)
	if err != nil {
		logp.Critical("processing events failed because of: %s", err)
		return err
	}
	bt.client.PublishEvents(processedEvents, publisher.Guaranteed)
	logp.Debug("api", "audit events sent: %v", processedEvents)
	return nil
}
