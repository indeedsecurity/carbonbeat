package app

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/indeedsecurity/carbonbeat/carbonclient"
)

func (bt *Carbonbeat) processNotifications(n carbonclient.Notifications) ([]common.MapStr, error) {
	var events []common.MapStr
	if n.Success {
		logp.Debug("api", "%v events collected", len(n.Events))

		for _, e := range n.Events {
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"timestamp":  e.EventTime,
				"type":       bt.config.Type,
				"url":        e.URL,
				"src_ip":     e.DeviceInfo.InternalIPAddress,
				"src_host":   e.DeviceInfo.DeviceName,
				"dst_ip":     e.DeviceInfo.ExternalIPAddress,
				"user":       e.DeviceInfo.Email,
				"cb": common.MapStr{
					"notificationType":   e.Type,
					"ruleName":           e.RuleName,
					"indicators":         e.ThreatInfo.Indicators,
					"incidentScore":      e.ThreatInfo.Score,
					"summary":            e.ThreatInfo.Summary,
					"deviceVersion":      e.DeviceInfo.DeviceVersion,
					"deviceType":         e.DeviceInfo.DeviceType,
					"policyName":         e.DeviceInfo.GroupName,
					"targetPriorityType": e.DeviceInfo.TargetPriorityType,
					"targetPriorityCode": e.DeviceInfo.TargetPriorityCode,
				},
			}
			events = append(events, event)
		}
		return events, nil
	}
	logp.Warn("something went wrong the thing, because notifications['success'] was false for what ever reason. good luck."+
		"here's whatever they gave us: %v", n)
	return events, nil
}
