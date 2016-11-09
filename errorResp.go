package box

import "errors"

// Box error messages
var (
	ErrItemNameInUse = errors.New("Item with the same name already exists")
)

// ErrorCodeResponse is a Box.net error code response.
// See https://docs.box.com/reference#section-error-code-response
type ErrorCodeResponse struct {
	Type        string                              `json:"type"`
	Status      int                                 `json:"status"`
	Code        string                              `json:"code"`
	ContextInfo map[string][]map[string]interface{} `json:"context_info"`
	HelpURL     string                              `json:"help_url"`
	Message     string                              `json:"message"`
	RequestID   string                              `json:"request_id"`
}

// ErrorContextInfo is an array of Error structs
// See https://docs.box.com/reference#section-error-code-response
type ErrorContextInfo struct {
	Errors []Error `json:"errors"`
}

// Error is error object which gives context about the error.
// See https://docs.box.com/reference#section-error-code-response
type Error struct {
	Reason  string `json:"reason"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
