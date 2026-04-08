package response

// Response is the unified struct for all API responses (success and error)
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewResponse creates a new Response instance with the given status, message, and result
func NewResponse(status string, message string, result interface{}) *Response {
	return &Response{
		Code:    status,
		Message: message,
		Data:    result,
	}
}
