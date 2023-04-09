package server

type baseResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
