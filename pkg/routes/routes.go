package routes

import (
	"noversystem/pkg/controller"
	"noversystem/pkg/dao"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	app.Use(logger.New())

	api := app.Group("/api")
	
    // Inisialisasi DAO dan Controller
	userDAO := dao.NewUserDao(db)
	authController := controller.NewAuthController(userDAO)

	// Grup rute untuk otentikasi
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
}