package api

type CommonError struct {
	Error string `json:"error"`
}

type CommonResponse struct {
	Message string `json:"message"`
}
