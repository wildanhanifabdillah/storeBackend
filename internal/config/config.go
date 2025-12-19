package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	// Database
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort string

	// SMTP
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string

	// Midtrans
	MidtransServerKey string
	MidtransClientKey string
}

func Load() *Config {
	// Load .env (ignore error for production)
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),

		DBHost: os.Getenv("DB_HOST"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		DBPort: os.Getenv("DB_PORT"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPEmail:    os.Getenv("SMTP_EMAIL"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),

		MidtransServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
		MidtransClientKey: os.Getenv("MIDTRANS_CLIENT_KEY"),
	}

	cfg.validate()
	return cfg
}

func (c *Config) validate() {
	required := map[string]string{
		"DB_HOST":             c.DBHost,
		"DB_USER":             c.DBUser,
		"DB_NAME":             c.DBName,
		"DB_PORT":             c.DBPort,
		"SMTP_HOST":           c.SMTPHost,
		"SMTP_EMAIL":          c.SMTPEmail,
		"SMTP_PASSWORD":       c.SMTPPassword,
		"MIDTRANS_SERVER_KEY": c.MidtransServerKey,
		"MIDTRANS_CLIENT_KEY": c.MidtransClientKey,
	}

	for k, v := range required {
		if v == "" {
			log.Fatalf("❌ Missing required env: %s", k)
		}
	}
}
