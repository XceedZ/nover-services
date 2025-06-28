package routes

import (
	"noversystem/pkg/controllers"
	"noversystem/pkg/dao"
	"noversystem/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	app.Use(logger.New())

	api := app.Group("/api")

	userDAO := dao.NewUserDao(db)

	// --- Auth Routes ---
	authController := controllers.NewAuthController(userDAO)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	// --- API v1 Group ---
	apiV1 := api.Group("/v1")

	// --- Bank Routes ---
	bankController := controllers.NewBankController(db)
	bankGroup := apiV1.Group("/bank")
	bankGroup.Get("/get", bankController.GetBankList)

	// --- User Routes ---
	userController := controllers.NewUserController(userDAO)
	userGroup := apiV1.Group("/user")

	// Protected endpoints
	protectedUserGroup := userGroup.Group("/", middleware.Protected())
	protectedUserGroup.Post("/request-author", userController.RequestBecomeAuthor)
	protectedUserGroup.Get("/author-status", userController.CheckAuthorStatus)
}
