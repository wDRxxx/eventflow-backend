package models

type OrderMessage struct {
	To []string

	TicketID    string
	EventTitle  string
	ImageURL    string
	RedirectURL string
}
