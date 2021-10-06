package internal

type (
	// Response contains details to be sent as a response to tigo
	// when tigo make callback request, name check request or payment
	// request.
	Response struct {
		StatusCode int               `json:"status_code"`
		Payload    interface{}       `json:"payload"`
		Headers    map[string]string `json:"headers"`
		Error      error             `json:"error"`
	}

	ResponseOption func(response *Response)
)

// NewResponse create a response to be sent back to Tigo. HTTP Status code, payload and its
// type need to be specified. Other fields like  Response.Error and Response.Headers can be
// changed using WithMoreResponseHeaders (add headers), WithResponseHeaders (replace all the
// existing ) and WithResponseError to add error its default value is nil, default value of
// Response.Headers is
// defaultResponseHeader = map[string]string{
//		"Content-Type": ContentTypeXml,
// }
func NewResponse(status int, payload interface{}, opts ...ResponseOption) *Response {

	var (
		defaultResponseHeader = map[string]string{
			"Content-Type": cTypeJson,
		}
	)

	response := &Response{
		StatusCode: status,
		Payload:    payload,
		Headers:    defaultResponseHeader,
		Error:      nil,
	}

	for _, opt := range opts {
		opt(response)
	}

	return response
}

func WithResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *Response) {
		response.Headers = headers
	}
}

func WithMoreResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *Response) {
		for key, value := range headers {
			response.Headers[key] = value
		}
	}
}

func WithDefaultXMLHeader() ResponseOption {
	return func(response *Response) {
		response.Headers["Content-Type"] = cTypeAppXml
	}
}

func WithResponseError(err error) ResponseOption {
	return func(response *Response) {
		response.Error = err
	}
}
