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

// BookCommentController menangani logika HTTP untuk komentar pada buku.
type BookCommentController struct {
	commentDAO *dao.BookCommentDao
	bookDAO    *dao.BookDao
	userDAO    *dao.UserDao
	log        *logrus.Logger
}

// NewBookCommentController membuat instance baru dari BookCommentController.
func NewBookCommentController(commentDAO *dao.BookCommentDao, bookDAO *dao.BookDao, userDAO *dao.UserDao) *BookCommentController {
	return &BookCommentController{
		commentDAO: commentDAO,
		bookDAO:    bookDAO,
		userDAO:    userDAO,
		log:        logrus.New(),
	}
}

type BookCommentListResponse struct {
	Comments []tables.BookComment `json:"comments"`
}

// CreateCommentRequest adalah DTO (Data Transfer Object) untuk request pembuatan komentar.
type CreateCommentRequest struct {
	CommentText     string `json:"commentText"`
	ParentCommentID *int64 `json:"parentCommentId,omitempty"`
}

// CreateBookComment adalah handler untuk membuat komentar pada sebuah buku.
// @Summary      Buat Komentar Buku
// @Description  Membuat komentar baru (atau balasan) pada sebuah buku oleh pengguna terotentikasi.
// @Tags         Book
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookId path int true "ID dari buku yang dikomentari"
// @Param        comment_data body CreateCommentRequest true "Isi komentar"
// @Success      201 {object} tables.BookComment
// @Failure      400 {object} ErrorResponse "Input tidak valid"
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      404 {object} ErrorResponse "Buku tidak ditemukan"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/books/{bookId}/comments [POST]
func (c *BookCommentController) CreateBookComment(ctx *fiber.Ctx) error {
	// 1. Dapatkan ID pengguna (aktor) dari token
	actorId, ok := ctx.Locals("userId").(int64)
	if !ok || actorId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid access."})
	}

	// 2. Dapatkan ID buku dari URL
	bookId, err := strconv.ParseInt(ctx.Params("bookId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid book ID."})
	}

	// 3. Parse request body
	var payload CreateCommentRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Cannot parse request body."})
	}

	// 4. Validasi input
	payload.CommentText = strings.TrimSpace(payload.CommentText)
	if payload.CommentText == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeAuthInputRequired, Message: "Comment text cannot be empty."})
	}

	// 5. Dapatkan data tambahan untuk notifikasi
	actor, err := c.userDAO.FindUserByID(ctx.Context(), actorId)
	if err != nil || actor == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to verify commenting user."})
	}

	book, err := c.bookDAO.GetBookDetailByID(ctx.Context(), bookId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get book details for notification."})
	}
	if book == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: constants.ErrCodeBookNotFound, Message: "Book not found."})
	}
	
	// Siapkan nama untuk notifikasi dengan fallback
	var actorName string
	if actor.PenName != nil && *actor.PenName != "" {
		actorName = *actor.PenName
	} else {
		actorName = "Seseorang"
	}
	
	// 6. Siapkan data dan panggil DAO
	commentData := &tables.BookComment{
		BookID:          bookId,
		UserID:          actorId,
		CommentText:     payload.CommentText,
		ParentCommentID: payload.ParentCommentID,
	}

	createdComment, err := c.commentDAO.CreateCommentAndNotify(ctx.Context(), commentData, actorName, book.Title)
	if err != nil {
		c.log.WithError(err).Error("Gagal membuat komentar buku di DAO")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to post comment."})
	}

	// 7. Lengkapi response dengan data aktor
	createdComment.AuthorPenName = actor.PenName
	if actor.AvatarURL != "" {
		createdComment.AuthorAvatar = &actor.AvatarURL
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdComment)
}

// GetBookComments adalah handler untuk mengambil semua komentar pada sebuah buku.
// @Summary      Dapatkan Komentar Buku
// @Description  Mengambil daftar semua komentar pada sebuah buku.
// @Tags         Book
// @Produce      json
// @Param        bookId path int true "ID dari buku"
// @Success      200 {object} BookCommentListResponse
// @Failure      400 {object} ErrorResponse "Input tidak valid"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/books/{bookId}/comments [GET]
func (c *BookCommentController) GetBookComments(ctx *fiber.Ctx) error {
	bookId, err := strconv.ParseInt(ctx.Params("bookId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid book ID."})
	}

	comments, err := c.commentDAO.GetCommentsByBookID(ctx.Context(), bookId)
	if err != nil {
		c.log.WithError(err).Error("Gagal mengambil komentar buku dari DAO")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve comments."})
	}

	// Memastikan response adalah array kosong `[]` bukan `null` jika tidak ada komentar
	if comments == nil {
		comments = make([]tables.BookComment, 0)
	}

	return ctx.Status(fiber.StatusOK).JSON(BookCommentListResponse{
		Comments: comments,
	})
}