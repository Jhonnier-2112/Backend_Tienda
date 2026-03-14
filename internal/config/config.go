package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	ServerAddress        string
	MercadoPagoAccessTok string
	PayPalClientID       string
	PayPalSecret         string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env not found")
	}
	return &Config{
		Port:                 getEnv("PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "Colombia1."),
		DBName:               getEnv("DB_NAME", "virtual_store"),
		JWTSecret:            getEnv("JWT_SECRET", "super_secret_key"),
		ServerAddress:        getEnv("VIRTUAL_STORE_SERVER_ADDRESS", ":10000"),
		MercadoPagoAccessTok: getEnv("MERCADOPAGO_ACCESS_TOKEN", ""),
		PayPalClientID:       getEnv("PAYPAL_CLIENT_ID", ""),
		PayPalSecret:         getEnv("PAYPAL_SECRET", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
