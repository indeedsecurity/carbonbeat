package app

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/indeedsecurity/carbonbeat/v2/carbonclient"
)

func (bt *Carbonbeat) processAuditEvents(ae carbonclient.AuditEvents) ([]common.MapStr, error) {
	var events []common.MapStr
	if ae.Success {
		logp.Debug("api", "%v events collected", len(ae.AuditEvents))

		for _, e := range ae.AuditEvents {
			event := common.MapStr{
				"@timestamp":  common.Time(time.Now()),
				"timestamp":   e.EventTime,
				"cb_type":     "audit",
				"eventId":     e.EventID,
				"loginName":   e.LoginName,
				"orgName":     e.OrgName,
				"flagged":     e.Flagged,
				"clientIp":    e.ClientIP,
				"verbose":     e.Verbose,
				"description": e.Description,
			}
			events = append(events, event)
		}
		return events, nil
	}
	logp.Warn("something went wrong, because notifications['success'] was false for what ever reason. good luck."+
		"here's whatever they gave us: %v", ae)
	return events, nil
}
