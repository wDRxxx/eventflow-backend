package models

import "time"

type User struct {
	ID         int64  `json:"-" db:"id"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"-" db:"password"`
	TGUsername string `json:"tg_username" db:"tg_username"`

	YookassaSettings *YookassaSettings `json:"yookassa_settings" db:"-"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type Event struct {
	ID            int64     `json:"-" db:"id"`
	Title         string    `json:"title" db:"title"`
	URLTitle      string    `json:"url_title,omitempty" db:"url_title"`
	Description   string    `json:"description" db:"description"`
	BeginningTime time.Time `json:"beginning_time" db:"beginning_time"`
	EndTime       time.Time `json:"end_time" db:"end_time"`
	CreatorID     int64     `json:"creator_id,omitempty" db:"creator_id"`
	IsPublic      bool      `json:"-" db:"is_public"`
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
	ID       int64  `json:"-" db:"id"`
	Price    int64  `json:"price" db:"price"`
	Currency string `json:"currency" db:"currency"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type YookassaSettings struct {
	ID      int64  `json:"-" db:"id"`
	UserID  int64  `json:"-" db:"user_id"`
	ShopID  int64  `json:"shop_id" db:"shop_id"`
	ShopKey string `json:"shop_key" db:"shop_key"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type Ticket struct {
	ID        string `json:"-" db:"id"`
	UserID    int64  `json:"-" db:"user_id"`
	EventID   int64  `json:"-" db:"event_id"`
	IsUsed    bool   `json:"is_used" db:"is_used"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
}
