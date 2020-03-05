package carbonclient

//AuditEvent is the structure of one Audit Log Event from the Carbonblack API
type AuditEvent struct {
	RequestURL  interface{} `json:"requestUrl"`
	EventTime   int64       `json:"eventTime"`
	EventID     string      `json:"eventId"`
	LoginName   string      `json:"loginName"`
	OrgName     string      `json:"orgName"`
	Flagged     bool        `json:"flagged"`
	ClientIP    string      `json:"clientIp"`
	Verbose     bool        `json:"verbose"`
	Description string      `json:"description"`
}

//AuditEvents is the array of audit log events sent from the Carbonblack API
type AuditEvents struct {
	AuditEvents []AuditEvent `json:"notifications"`
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
}
