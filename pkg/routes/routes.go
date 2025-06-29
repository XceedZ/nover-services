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
	chapterDAO := dao.NewChapterDao(db) // Inisialisasi ChapterDAO
	reviewDAO := dao.NewReviewDao(db)

	// --- Auth Routes ---
	authController := controllers.NewAuthController(userDAO)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	// --- API v1 Group ---
	apiV1 := api.Group("/v1")

	// --- Genre Routes (Public) ---
	genreController := controllers.NewGenreController(genreDAO)
	apiV1.Get("/genres", genreController.GetAllGenres)

	// --- Bank Routes (Public) ---
	bankController := controllers.NewBankController(db)
	bankGroup := apiV1.Group("/bank")
	bankGroup.Get("/get", bankController.GetBankList)

	// --- User Routes (Protected) ---
	userController := controllers.NewUserController(userDAO)
	userGroup := apiV1.Group("/user")
	protectedUserGroup := userGroup.Group("/", middleware.Protected())
	protectedUserGroup.Post("/request-author", userController.RequestBecomeAuthor)
	protectedUserGroup.Get("/author-status", userController.CheckAuthorStatus)

	// --- Book Routes ---
	bookController := controllers.NewBookController(bookDAO, userDAO, chapterDAO, reviewDAO)

	// 👉 PUBLIC Book Endpoints (tidak pakai middleware, bebas akses tanpa token)
	apiV1.Get("/books", bookController.GetPublishedBookList)
	apiV1.Get("/authors/:authorId/books", bookController.GetBooksByAuthor)
	apiV1.Get("/chapters/:chapterId", controllers.NewChapterController(chapterDAO, bookDAO).GetChapterContent)

	// 👉 PROTECTED Book Endpoints (wajib pakai token)
	bookGroup := apiV1.Group("/books", middleware.Protected())
	bookGroup.Post("/create", bookController.CreateBook)
	bookGroup.Get("/my-books", bookController.GetMyBooks)
	bookGroup.Patch("/:bookId/publish", bookController.PublishBook)
	bookGroup.Patch("/:bookId/unpublish", bookController.UnpublishBook)
	bookGroup.Patch("/:bookId/complete", bookController.CompleteBook)
	bookGroup.Patch("/:bookId/hold", bookController.HoldBook)
	bookGroup.Get("/:bookId/detail", bookController.GetMyBookDetail)

	// Chapter creation (Protected, karena di bawah bookGroup)
	chapterController := controllers.NewChapterController(chapterDAO, bookDAO)
	bookGroup.Post("/:bookId/chapters", chapterController.CreateChapter)
	apiV1.Get("/books/:bookId", bookController.GetPublicBookDetail)
}
