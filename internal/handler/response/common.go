package response

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(status string, message string, result interface{}) *Response {
	return &Response{status, message, result}
}
