package msg

// IDMessage is a prototype of a generic json message containing element's id
type IDMessage struct {
	RequestMessage
	ID string `json:"id"`
}

// RequestMessage is a generic request message
type RequestMessage struct {
	RequestID string `json:"requestID"`
	CallerID  string `json:"callerID,omitempty"`
}

// ResponseMessage is a generic acknowledge response containing http-like status code and error details if necessary
type ResponseMessage struct {
	RequestID string  `json:"requestID"`
	CallerID  string  `json:"callerID,omitempty"`
	Status    int     `json:"status,omitempty"`
	Error     *string `json:"error,omitempty"`
}
