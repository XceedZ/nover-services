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

// ChapterController menangani logika HTTP yang terkait dengan chapter.
type ChapterController struct {
	chapterDAO *dao.ChapterDao
	bookDAO    *dao.BookDao // Diperlukan untuk validasi kepemilikan buku
	log        *logrus.Logger
}

// NewChapterController membuat instance baru dari ChapterController.
func NewChapterController(chapterDAO *dao.ChapterDao, bookDAO *dao.BookDao) *ChapterController {
	return &ChapterController{
		chapterDAO: chapterDAO,
		bookDAO:    bookDAO,
		log:        logrus.New(),
	}
}

// CreateChapterRequest adalah payload untuk menambah chapter baru.
type CreateChapterRequest struct {
	Title    string  `json:"title"`
	Content  *string `json:"content"`
	CoinCost int     `json:"coinCost"`
}

// CreateChapter adalah handler untuk menambah chapter baru ke sebuah buku.
// @Summary      Tambah Chapter Baru
// @Description  Menambahkan sebuah chapter baru ke buku yang sudah ada. Pengguna harus menjadi pemilik buku tersebut.
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookId path int true "ID Buku"
// @Param        chapter_data body CreateChapterRequest true "Data chapter yang akan dibuat"
// @Success      201 {object} tables.Chapter
// @Failure      400 {object} ErrorResponse "Input tidak valid"
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      403 {object} ErrorResponse "Akses ditolak (bukan pemilik buku)"
// @Failure      404 {object} ErrorResponse "Buku tidak ditemukan"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/books/{bookId}/chapters [POST]
func (c *ChapterController) CreateChapter(ctx *fiber.Ctx) error {
	// 1. Ambil userID dari token JWT
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid user token."})
	}

	// 2. Ambil bookId dari parameter URL
	bookId, err := strconv.ParseInt(ctx.Params("bookId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid book ID."})
	}

	// 3. Validasi Kepemilikan Buku
	book, err := c.bookDAO.GetBookWithAuthor(ctx.Context(), bookId)
	if err != nil {
		c.log.WithError(err).Error("Gagal mengambil detail buku untuk validasi kepemilikan")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get book details."})
	}
	if book == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: constants.ErrCodeUserNotFound, Message: "Book not found."})
	}
	if book.AuthorID != userId {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Code: constants.ErrCodeBookNotOwner, Message: "You are not the owner of this book."})
	}

	// 4. Parse payload request
	var payload CreateChapterRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Cannot parse request body."})
	}

	payload.Title = strings.TrimSpace(payload.Title)
	if payload.Title == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeAuthInputRequired, Message: "Chapter title is required."})
	}

	// 5. Tentukan urutan chapter baru
	lastOrder, err := c.chapterDAO.GetLastChapterOrder(ctx.Context(), bookId)
	if err != nil {
		c.log.WithError(err).Error("Gagal mendapatkan urutan chapter terakhir")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to determine chapter order."})
	}
	newOrder := lastOrder + 1

	// 6. Siapkan data dan panggil DAO untuk membuat chapter
	newChapter := &tables.Chapter{
		BookID:       bookId,
		Title:        payload.Title,
		Content:      payload.Content,
		ChapterOrder: newOrder,
		CoinCost:     payload.CoinCost,
	}

	createdChapter, err := c.chapterDAO.CreateChapter(ctx.Context(), newChapter)
	if err != nil {
		c.log.WithError(err).Error("Gagal membuat chapter baru di DAO")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to create new chapter."})
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdChapter)
}
