package controllers

import (
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"
	"noversystem/pkg/tables"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// BookController menangani logika HTTP yang terkait dengan buku.
type BookController struct {
	bookDAO *dao.BookDao
	userDAO *dao.UserDao
	log     *logrus.Logger
}

// NewBookController membuat instance baru dari BookController.
func NewBookController(bookDAO *dao.BookDao, userDAO *dao.UserDao) *BookController {
	return &BookController{
		bookDAO: bookDAO,
		userDAO: userDAO,
		log:     logrus.New(),
	}
}

// CreateBookRequest adalah struktur payload untuk membuat buku baru.
type CreateBookRequest struct {
	Title         string  `json:"title"`
	Description   *string `json:"description"`
	CoverImageURL *string `json:"coverImageUrl"`
	GenreIDs      []int64 `json:"genreIds"`
}

// --- STRUCT BARU UNTUK RESPONSE ---
// BookListResponse adalah struktur untuk response daftar buku yang dibungkus.
type BookListResponse struct {
	BookList []tables.Book `json:"bookList"`
}

// CreateBook adalah handler untuk endpoint pembuatan buku baru.
// @Summary      Buat Buku Baru
// @Description  Membuat buku baru oleh pengguna yang sudah terotentikasi dan berstatus sebagai penulis.
// @Tags         Book
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        book_data body CreateBookRequest true "Data buku yang akan dibuat"
// @Success      201 {object} tables.Book
// @Failure      403 {object} ErrorResponse "Akses ditolak (bukan penulis)"
// @Router       /v1/books/create [POST]
func (c *BookController) CreateBook(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid access."})
	}
	author, err := c.userDAO.FindUserByID(ctx.Context(), userId)
	if err != nil || author == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to verify author status."})
	}
	if author.FlgAuthor != "Y" {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Code: "auth.error.not_an_author", Message: "Access denied. Only authors can create books."})
	}
	var payload CreateBookRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Cannot parse request body."})
	}
	payload.Title = strings.TrimSpace(payload.Title)
	if payload.Title == "" || len(payload.GenreIDs) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeAuthInputRequired, Message: "Title and at least one genre ID are required."})
	}
	bookData := &tables.Book{
		Title:         payload.Title,
		Description:   payload.Description,
		CoverImageURL: payload.CoverImageURL,
	}
	createdBook, err := c.bookDAO.CreateBook(ctx.Context(), bookData, userId, payload.GenreIDs)
	if err != nil {
		c.log.WithError(err).Error("Gagal membuat buku baru di DAO")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to create new book."})
	}
	return ctx.Status(fiber.StatusCreated).JSON(createdBook)
}

// GetMyBooks adalah handler untuk mendapatkan daftar buku milik penulis yang sedang login (dilindungi).
// @Summary      Dapatkan Buku Saya (Pribadi)
// @Description  Mengambil daftar semua buku yang ditulis oleh pengguna yang sedang login.
// @Tags         Book
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} BookListResponse
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/books/my-books [GET]
func (c *BookController) GetMyBooks(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid access, user not authenticated properly."})
	}

	books, err := c.bookDAO.GetBooksByAuthorID(ctx.Context(), userId)
	if err != nil {
		c.log.WithError(err).Errorf("Gagal mengambil buku untuk author ID %d", userId)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve books."})
	}
	
	// Jika tidak ada buku, kembalikan array kosong di dalam objek.
	if books == nil {
		books = []tables.Book{}
	}

	// Bungkus slice 'books' di dalam struct BookListResponse
	response := BookListResponse{BookList: books}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetBooksByAuthor adalah handler publik untuk mendapatkan daftar buku dari seorang penulis.
// @Summary      Dapatkan Buku Berdasarkan Penulis (Publik)
// @Description  Mengambil daftar semua buku dari seorang penulis berdasarkan ID penulis.
// @Tags         Book
// @Produce      json
// @Param        authorId path int true "ID dari Penulis"
// @Success      200 {object} BookListResponse
// @Failure      400 {object} ErrorResponse "ID Penulis tidak valid"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/authors/{authorId}/books [GET]
func (c *BookController) GetBooksByAuthor(ctx *fiber.Ctx) error {
	authorIdParam := ctx.Params("authorId")
	authorId, err := strconv.ParseInt(authorIdParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid author ID format."})
	}

	books, err := c.bookDAO.GetBooksByAuthorID(ctx.Context(), authorId)
	if err != nil {
		c.log.WithError(err).Errorf("Gagal mengambil buku untuk author ID %d", authorId)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve author's books."})
	}

	if books == nil {
		books = []tables.Book{}
	}

	// Bungkus slice 'books' di dalam struct BookListResponse
	response := BookListResponse{BookList: books}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
