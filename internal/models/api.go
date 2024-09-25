package models

type DefaultResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
