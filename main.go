package main

import (
	"context"

	_ "noversystem/docs" // PENTING: Import blank direktori docs yang akan kita generate

	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"

	"noversystem/pkg/config" // Import package config yang kita buat
	"noversystem/pkg/routes" // Import package routes yang kita buat
)

// @title Nover System API
// @version 1.0
// @description Ini adalah dokumentasi API untuk Nover System.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@noversystem.dev
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
func main() {
	// 1. Muat Konfigurasi
	if err := config.LoadConfig(); err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}
	cfg := config.Cfg

	// 2. Inisialisasi Logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(cfg.Log.Level)
	logrus.Infof("Starting %s application...", cfg.App.Name)

	// 3. Inisialisasi Koneksi Database (PostgreSQL dengan pgx)
		config, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		logrus.Fatalf("Unable to parse database configuration: %v", err)
	}

	// 2. Atur Query Exec Mode untuk menonaktifkan prepared statement cache
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// 3. Hubungkan ke database menggunakan config yang sudah diubah
    // Perhatikan: di v5 menggunakan pgxpool.NewWithConfig
	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logrus.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()
	logrus.Info("Successfully connected to the database using pgx/v5")

	// 4. Inisialisasi Fiber App
	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})

	app.Get("/api/docs/*", swagger.HandlerDefault)

	routes.SetupRoutes(app, db)

	// Gunakan middleware recover agar aplikasi tidak crash jika terjadi panic
	app.Use(recover.New())

	// 5. Setup Rute API
	routes.SetupRoutes(app, db)
	logrus.Info("API routes have been initialized")

	// 6. Jalankan Server
	listenAddr := ":" + cfg.App.Port
	logrus.Infof("Server is starting and listening on port %s", cfg.App.Port)

	if err := app.Listen(listenAddr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
