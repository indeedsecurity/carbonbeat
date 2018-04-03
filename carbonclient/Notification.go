package carbonclient

// Event is the structure of one event from CB
type Event struct {
	ThreatInfo struct {
		IncidentID string `json:"incidentId"`
		Score      int    `json:"score"`
		Summary    string `json:"summary"`
		Indicators []struct {
			ApplicationName string `json:"applicationName"`
			Sha256Hash      string `json:"sha256Hash"`
			IndicatorName   string `json:"indicatorName"`
		} `json:"indicators"`
		Time int64 `json:"time"`
	} `json:"threatInfo"`
	URL              string `json:"url"`
	EventTime        int64  `json:"eventTime"`
	EventID          string `json:"eventId"`
	EventDescription string `json:"eventDescription"`
	DeviceInfo       struct {
		ExternalIPAddress  string      `json:"externalIpAddress"`
		DeviceName         string      `json:"deviceName"`
		DeviceHostName     interface{} `json:"deviceHostName"`
		DeviceVersion      string      `json:"deviceVersion"`
		DeviceID           int         `json:"deviceId"`
		Email              string      `json:"email"`
		GroupName          string      `json:"groupName"`
		InternalIPAddress  string      `json:"internalIpAddress"`
		DeviceType         string      `json:"deviceType"`
		TargetPriorityType string      `json:"targetPriorityType"`
		TargetPriorityCode int         `json:"targetPriorityCode"`
	} `json:"deviceInfo"`
	RuleName string `json:"ruleName"`
	Type     string `json:"type"`
}

// Notifications is the array of Events sent back from CB
type Notifications struct {
	Events  []Event `json:"notifications"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
}
