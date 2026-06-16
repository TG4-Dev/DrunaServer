package model

type ErrorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

type APIResponse struct {
	Data  interface{} `json:"data"`
	Error *ErrorBody  `json:"error"`
}
