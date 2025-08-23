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
    bookCommentDAO := dao.NewBookCommentDao(db) // ✨ 1. Inisialisasi DAO baru
    notificationDAO := dao.NewNotificationDao(db) // ✨ Inisialisasi DAO baru
    walletDAO := dao.NewWalletDao(db) // ✨ 1. Inisialisasi WalletDAO
	transactionDAO := dao.NewTransactionDao(db) // ✨ 1. Inisialisasi TransactionDAO
	checkinDAO := dao.NewCheckinDao(db)    // ✨ Inisialisasi DAO baru
	missionDAO := dao.NewMissionDao(db)    // ✨ Inisialisasi DAO baru

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
    bookCommentController := controllers.NewBookCommentController(bookCommentDAO, bookDAO, userDAO)
    notificationController := controllers.NewNotificationController(notificationDAO) // ✨ Inisialisasi Controller baru
    walletController := controllers.NewWalletController(walletDAO) // ✨ 2. Inisialisasi WalletController
	transactionController := controllers.NewTransactionController(transactionDAO) // ✨ 3. Inisialisasi TransactionController
	checkinController := controllers.NewCheckinController(checkinDAO) // ✨ Inisialisasi Controller baru
	missionController := controllers.NewMissionController(missionDAO) // ✨ Inisialisasi Controller baru

	// 👉 PUBLIC Book Endpoints (tidak pakai middleware, bebas akses tanpa token)
	apiV1.Get("/books", bookController.GetPublishedBookList)
	apiV1.Get("/authors/:authorId/books", bookController.GetBooksByAuthor)
	apiV1.Get("/chapters/:chapterId", controllers.NewChapterController(chapterDAO, bookDAO).GetChapterContent)

	apiV1.Get("/books/:bookId/comments", bookCommentController.GetBookComments)

	// 👉 PROTECTED Book Endpoints (wajib pakai token)
	bookGroup := apiV1.Group("/books", middleware.Protected())
	bookGroup.Post("/create", bookController.CreateBook)
	bookGroup.Get("/my-books", bookController.GetMyBooks)
	bookGroup.Patch("/:bookId/publish", bookController.PublishBook)
	bookGroup.Patch("/:bookId/unpublish", bookController.UnpublishBook)
	bookGroup.Patch("/:bookId/complete", bookController.CompleteBook)
	bookGroup.Patch("/:bookId/hold", bookController.HoldBook)
	bookGroup.Get("/:bookId/detail", bookController.GetMyBookDetail)
    bookGroup.Post("/:bookId/comments", bookCommentController.CreateBookComment)

	// Chapter creation (Protected, karena di bawah bookGroup)
	chapterController := controllers.NewChapterController(chapterDAO, bookDAO)
	bookGroup.Post("/:bookId/chapters", chapterController.CreateChapter)
	apiV1.Get("/books/:bookId", bookController.GetPublicBookDetail)

	notifGroup := apiV1.Group("/notifications", middleware.Protected())
    notifGroup.Get("/", notificationController.GetNotifications) // ✨ Daftarkan Route GET baru

	walletGroup := apiV1.Group("/wallet", middleware.Protected())
    walletGroup.Get("/my-balance", walletController.GetMyWallet)
    walletGroup.Get("/transactions", transactionController.GetMyTransactions)

	eventGroup := apiV1.Group("/events", middleware.Protected())
	eventGroup.Get("/check-in/status", checkinController.GetStatus)
	eventGroup.Post("/check-in", checkinController.CheckIn)
	eventGroup.Get("/missions/daily", missionController.GetDailyMissions)
}
