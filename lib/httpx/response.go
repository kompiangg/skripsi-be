package httpx

type HTTPBaseResponse struct {
	Error *HTTPErrorBaseResponse `json:"error"`
	Data  interface{}            `json:"data"`
}

type HTTPErrorBaseResponse struct {
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
}
