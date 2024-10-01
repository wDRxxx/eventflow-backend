package mailer

import (
	"github.com/wDRxxx/eventflow-backend/internal/models"
)

type Mailer interface {
	ListenForMails()
	SendHTMLMessage(body []byte, to []string) error
	SendOrderMail(msg *models.OrderMessage)
}
