package models

type DefaultResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type BuyTicketRequest struct {
	EventUrlTitle string `json:"event_url_title"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	PriceID       int64  `json:"price_id,omitempty"`
	UserEmail     string `json:"-"`
}
