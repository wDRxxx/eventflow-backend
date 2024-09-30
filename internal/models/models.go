package models

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	yoopayment "github.com/wDRxxx/yookassa-go-sdk/yookassa/models/payment"
)

type User struct {
	ID         int64  `json:"-" db:"id"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password,omitempty" db:"password"`
	TGUsername string `json:"tg_username,omitempty" db:"tg_username"`

	YookassaSettings YookassaSettings `json:"yookassa_settings" db:"-"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type Event struct {
	ID            int64     `json:"-" db:"id"`
	Title         string    `json:"title" db:"title"`
	URLTitle      string    `json:"url_title,omitempty" db:"url_title"`
	Description   string    `json:"description" db:"description"`
	BeginningTime time.Time `json:"beginning_time" db:"beginning_time"`
	EndTime       time.Time `json:"end_time" db:"end_time"`
	CreatorID     int64     `json:"creator_id,omitempty" db:"creator_id"`
	IsPublic      bool      `json:"is_public" db:"is_public"`
	Location      string    `json:"location" db:"location"`
	IsFree        bool      `json:"is_free" db:"is_free"`
	PreviewImage  string    `json:"preview_image" db:"preview_image"`
	UTCOffset     int64     `json:"utc_offset" db:"utc_offset"`
	Capacity      int64     `json:"capacity" db:"capacity"`
	MinimalAge    int64     `json:"minimal_age" db:"minimal_age"`
	Prices        []*Price  `json:"prices" db:"-"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type Price struct {
	ID       int64  `json:"id,omitempty" db:"id"`
	Price    int64  `json:"price" db:"price"`
	Currency string `json:"currency" db:"currency"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type YookassaSettings struct {
	ID      int64  `json:"-" db:"id"`
	UserID  int64  `json:"-" db:"user_id"`
	ShopID  string `json:"shop_id" db:"shop_id"`
	ShopKey string `json:"shop_key" db:"shop_key"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type Ticket struct {
	ID        string `json:"id" db:"id"`
	UserID    int64  `json:"-" db:"user_id"`
	User      User   `json:"u-" db:"-"`
	EventID   int64  `json:"-" db:"event_id"`
	Event     Event  `json:"event,omitempty" db:"-"`
	IsUsed    bool   `json:"is_used" db:"is_used"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	PaymentID string `json:"-" db:"payment_id"`
}

type TicketPayment struct {
	BuyTicketRequest *BuyTicketRequest
	Payment          *yoopayment.Payment
	User             *User
	Event            *Event
	Ctx              context.Context
}
