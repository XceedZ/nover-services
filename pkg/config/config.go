package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config menampung semua konfigurasi aplikasi
type Config struct {
	App struct {
		Name string
		Port string
	}
	DB struct {
		DSN string // Data Source Name
	}
	Log struct {
		Level logrus.Level
	}
}

// Cfg adalah variabel global untuk menampung konfigurasi yang sudah di-load
var Cfg *Config

// LoadConfig memuat konfigurasi dari file .env dan menginisialisasi variabel Cfg
func LoadConfig() error {
	// Muat file .env dari root direktori
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Error loading .env file, using environment variables")
	}

	Cfg = new(Config)

	// Konfigurasi Aplikasi
	Cfg.App.Name = os.Getenv("APP_NAME")
	Cfg.App.Port = os.Getenv("APP_PORT")
	if Cfg.App.Port == "" {
		Cfg.App.Port = "8080" // Default port
	}

	// Konfigurasi Database
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	_, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %w", err)
	}

	Cfg.DB.DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPortStr, dbUser, dbPassword, dbName, dbSSLMode)

	// Konfigurasi Logger
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel // Default log level
	}
	Cfg.Log.Level = logLevel

	return nil
}