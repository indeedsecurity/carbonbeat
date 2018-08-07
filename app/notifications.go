package app

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/indeedsecurity/carbonbeat/carbonclient"
)

func (bt *Carbonbeat) processNotifications(n carbonclient.Notifications) ([]common.MapStr, error) {
	var notifications []common.MapStr
	if n.Success {
		logp.Debug("api", "%v events collected", len(n.Notifications))

		for _, e := range n.Notifications {
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"timestamp":  e.EventTime,
				"type":       e.Type,
				"url":        e.URL,
				"src_ip":     e.DeviceInfo.InternalIPAddress,
				"src_host":   e.DeviceInfo.DeviceName,
				"dst_ip":     e.DeviceInfo.ExternalIPAddress,
				"user":       e.DeviceInfo.Email,
				"message":    e.EventDescription,
				"cb": common.MapStr{
					"event_id":             e.EventID,
					"rule_name":            e.RuleName,
					"device_version":       e.DeviceInfo.DeviceVersion,
					"device_type":          e.DeviceInfo.DeviceType,
					"group_name":           e.DeviceInfo.GroupName,
					"target_priority_type": e.DeviceInfo.TargetPriorityType,
					"target_priority_code": e.DeviceInfo.TargetPriorityCode,
				},
			}

			if e.ThreatInfo.IncidentID != "" {
				event.Put("cb.threat_info", e.ThreatInfo)
			}

			if e.PolicyAction.Action != "" {
				event.Put("cb.policy_action", e.PolicyAction)
			}

			notifications = append(notifications, event)
		}
		return notifications, nil
	}
	logp.Warn("something went wrong, because notifications['success'] was false for what ever reason. good luck."+
		"here's whatever they gave us: %v", n)
	return notifications, nil
}
