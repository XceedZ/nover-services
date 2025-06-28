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

	// --- Inisialisasi semua DAO ---
	userDAO := dao.NewUserDao(db)
	bookDAO := dao.NewBookDao(db)
	genreDAO := dao.NewGenreDao(db)

	// --- Auth Routes ---
	authController := controllers.NewAuthController(userDAO)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	// --- API v1 Group ---
	apiV1 := api.Group("/v1")

	genreController := controllers.NewGenreController(genreDAO)
	apiV1.Get("/genres", genreController.GetAllGenres)

	// --- Bank Routes ---
	bankController := controllers.NewBankController(db)
	bankGroup := apiV1.Group("/bank")
	bankGroup.Get("/get", bankController.GetBankList)

	// --- User Routes ---
	userController := controllers.NewUserController(userDAO)
	userGroup := apiV1.Group("/user")
	protectedUserGroup := userGroup.Group("/", middleware.Protected())
	protectedUserGroup.Post("/request-author", userController.RequestBecomeAuthor)
	protectedUserGroup.Get("/author-status", userController.CheckAuthorStatus)

	// --- Book Routes ---
	bookController := controllers.NewBookController(bookDAO, userDAO)
	bookGroup := apiV1.Group("/books", middleware.Protected())
	bookGroup.Post("/create", bookController.CreateBook)
	bookGroup.Get("/my-books", bookController.GetMyBooks)
	apiV1.Get("/authors/:authorId/books", bookController.GetBooksByAuthor)

}
