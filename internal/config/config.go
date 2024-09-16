package config

import "github.com/joho/godotenv"

func Load(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil {
		return err
	}

	return nil
}
