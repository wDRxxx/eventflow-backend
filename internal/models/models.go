package models

import "time"

type User struct {
	ID         int64  `json:"-"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	TGUsername string `json:"tg_username"`

	YookassaSettings *YookassaSettings `json:"yookassa_settings"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Event struct {
	ID            int64     `json:"-"`
	Title         string    `json:"title"`
	URLTitle      string    `json:"url_title"`
	Description   string    `json:"description"`
	BeginningTime time.Time `json:"beginning_time"`
	EndTime       time.Time `json:"end_time"`
	CreatorID     int64     `json:"-"`
	IsPublic      bool      `json:"-"`
	Location      string    `json:"location"`
	IsFree        bool      `json:"is_free"`
	PreviewImage  string    `json:"preview_image"`
	UTCOffset     int64     `json:"utc_offset"`
	Capacity      int64     `json:"capacity"`
	MinimalAge    int64     `json:"minimal_age"`
	Prices        []*Price  `json:"prices"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Price struct {
	ID       int64  `json:"-"`
	Price    int64  `json:"price"`
	Currency string `json:"currency"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type YookassaSettings struct {
	ID      int64  `json:"-"`
	UserID  int64  `json:"-"`
	ShopID  int64  `json:"shop_id"`
	ShopKey string `json:"shop_key"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
