package errorLogMessages

type ErrorLogMessage struct {
	Tenant    string `json:"tenant"`
	EpochDay  int    `json:"epochDay"`
	Time      string `json:"time"`
	MessageID string `json:"messageId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Reason    string `json:"reason"`
}

// AllUsersResponse to form payload of an array of User structs
type ErrorLogMessageResponse struct {
	ErrorLogMessages []ErrorLogMessage `json:"errorMessages"`
}
