package config

import (
	"os"
)

type YookassaConfig struct {
	shopID  string
	shopKey string
}

func (c *YookassaConfig) ShopID() string {
	return c.shopID
}

func (c *YookassaConfig) ShopKey() string {
	return c.shopKey
}

func NewYookassaConfig() *YookassaConfig {
	shopID := os.Getenv("YOOKASSA_SHOP_ID")
	if shopID == "" {
		panic("YOOKASSA_SHOP_ID environment variable is empty")
	}
	shopKey := os.Getenv("YOOKASSA_SHOP_KEY")
	if shopKey == "" {
		panic("YOOKASSA_SHOP_KEY environment variable is empty")
	}

	return &YookassaConfig{
		shopID:  shopID,
		shopKey: shopKey,
	}
}
