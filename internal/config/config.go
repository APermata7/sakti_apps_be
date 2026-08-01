package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DatabaseURL string

	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	SupabaseJWTURL     string

	OfficeLat    float64
	OfficeLon    float64
	OfficeRadius int

	CloudinaryCloudName    string
	CloudinaryAPIKey       string
	CloudinaryAPISecret    string
	CloudinaryUploadFolder string

	TelegramBotToken string
	FrontendURL      string
	ResetPasswordURL string

	CORSAllowedOrigins string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	return &Config{
		AppName: getEnv("APP_NAME", "sakti_apps_be"),
		AppEnv:  getEnv("APP_ENV", "local"),
		AppPort: getEnv("APP_PORT", "8080"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseJWTURL:     getEnv("SUPABASE_JWT_URL", ""),

		OfficeLat:    getEnvAsFloat("OFFICE_LAT", -7.942777),
		OfficeLon:    getEnvAsFloat("OFFICE_LON", 112.641110),
		OfficeRadius: getEnvAsInt("OFFICE_RADIUS", 500),

		CloudinaryCloudName:    getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:       getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret:    getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryUploadFolder: getEnv("CLOUDINARY_UPLOAD_FOLDER", "sakti-apps"),

		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:65138/"),
		ResetPasswordURL: getEnv("RESET_PASSWORD_URL", "sakti://reset-password"),

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:54323"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}