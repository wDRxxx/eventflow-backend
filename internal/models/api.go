package models

type DefaultResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
