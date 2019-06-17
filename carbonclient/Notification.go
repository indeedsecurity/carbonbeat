package carbonclient

// Notification is the structure of one notification event from CB
type Notification struct {
	PolicyAction struct {
		Sha256Hash      string `json:"sha256Hash"`
		Action          string `json:"action"`
		Reputation      string `json:"reputation"`
		ApplicationName string `json:"applicationName"`
	} `json:"policyAction,omitempty"`
	ThreatInfo struct {
		IncidentID string `json:"incidentId"`
		Score      int    `json:"score"`
		Summary    string `json:"summary"`
		Time       int64  `json:"time"`
		Indicators []struct {
			ApplicationName string `json:"applicationName"`
			Sha256Hash      string `json:"sha256Hash"`
			IndicatorName   string `json:"indicatorName"`
		} `json:"indicators"`
		ThreatCause struct {
			OriginSourceType string `json:"originSourceType"`
			Reputation       string `json:"reputation"`
			Actor            string `json:"actor"`
			ActorName        string `json:"actorName"`
			Reason           string `json:"reason"`
			ActorType        string `json:"actorType"`
			ThreatCategory   string `json:"threatCategory"`
			ActorProcessPPid string `json:"actorProcessPPid"`
			CauseEventID     string `json:"causeEventId"`
		} `json:"threatCause"`
	} `json:"threatInfo,omitempty"`
	URL              string `json:"url"`
	EventTime        int64  `json:"eventTime"`
	EventDescription string `json:"eventDescription"`
	DeviceInfo       struct {
		ExternalIPAddress  string `json:"externalIpAddress"`
		DeviceName         string `json:"deviceName"`
		DeviceHostName     string `json:"deviceHostName"`
		DeviceVersion      string `json:"deviceVersion"`
		DeviceID           int    `json:"deviceId"`
		Email              string `json:"email"`
		GroupName          string `json:"groupName"`
		InternalIPAddress  string `json:"internalIpAddress"`
		DeviceType         string `json:"deviceType"`
		TargetPriorityType string `json:"targetPriorityType"`
		TargetPriorityCode int    `json:"targetPriorityCode"`
	} `json:"deviceInfo"`
	RuleName string `json:"ruleName"`
	Type     string `json:"type"`
}

// Notifications is the array of Notification sent back from CB
type Notifications struct {
	Notifications []Notification `json:"notifications"`
	Success       bool           `json:"success"`
	Message       string         `json:"message"`
}
