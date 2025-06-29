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
	bookDAO    *dao.BookDao
	userDAO    *dao.UserDao
	chapterDAO *dao.ChapterDao
	reviewDAO  *dao.ReviewDao
	log        *logrus.Logger
}

// NewBookController membuat instance baru dari BookController.
func NewBookController(bookDAO *dao.BookDao, userDAO *dao.UserDao, chapterDAO *dao.ChapterDao, reviewDAO *dao.ReviewDao) *BookController {
	return &BookController{
		bookDAO:    bookDAO,
		userDAO:    userDAO,
		chapterDAO: chapterDAO,
		reviewDAO:  reviewDAO,
		log:        logrus.New(),
	}
}

// CreateBookRequest adalah struktur payload untuk membuat buku baru.
type CreateBookRequest struct {
	Title         string  `json:"title" example:"Judul Buku Fantasi"`
	Description   *string `json:"description" example:"Deskripsi singkat tentang petualangan epik."`
	CoverImageURL *string `json:"coverImageUrl" example:"https://path.to/your/image.jpg"`
	GenreIDs      []int64 `json:"genreIds" example:"1,2"`
}

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
// @Failure      400 {object} ErrorResponse "Input tidak valid"
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      403 {object} ErrorResponse "Akses ditolak (bukan penulis)"
// @Failure      500 {object} ErrorResponse "Error internal server"
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
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Code: constants.ErrCodeBookNotOwner, Message: "Access denied. Only authors can create books."})
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

// GetMyBooks adalah handler untuk mendapatkan daftar buku milik penulis yang sedang login.
// @Summary      Dapatkan Buku Saya (Pribadi)
// @Description  Mengambil daftar semua buku (termasuk draft) yang ditulis oleh pengguna yang sedang login.
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
	books, err := c.bookDAO.GetBooksByAuthorID(ctx.Context(), userId, false)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve your books."})
	}
	if books == nil {
		books = []tables.Book{}
	}
	return ctx.Status(fiber.StatusOK).JSON(BookListResponse{BookList: books})
}

// GetBooksByAuthor adalah handler publik untuk mendapatkan daftar buku dari seorang penulis.
// @Summary      Dapatkan Buku Berdasarkan Penulis (Publik)
// @Description  Mengambil daftar semua buku (kecuali draft) dari seorang penulis berdasarkan ID penulis.
// @Tags         Book
// @Produce      json
// @Param        authorId path int true "ID dari Penulis"
// @Success      200 {object} BookListResponse
// @Failure      400 {object} ErrorResponse "ID Penulis tidak valid"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/authors/{authorId}/books [GET]
func (c *BookController) GetBooksByAuthor(ctx *fiber.Ctx) error {
	authorId, err := strconv.ParseInt(ctx.Params("authorId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid author ID format."})
	}
	books, err := c.bookDAO.GetBooksByAuthorID(ctx.Context(), authorId, true)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve author's books."})
	}
	if books == nil {
		books = []tables.Book{}
	}
	return ctx.Status(fiber.StatusOK).JSON(BookListResponse{BookList: books})
}

// GetPublishedBookList adalah handler publik untuk mendapatkan daftar buku dengan pagination.
// @Summary      Dapatkan Daftar Buku (Publik, Paginasi)
// @Description  Mengambil daftar semua buku yang sudah dipublikasikan (status 'P', 'C', 'H') dengan sistem pagination.
// @Tags         Book
// @Produce      json
// @Param        page query int false "Nomor Halaman" default(1)
// @Param        limit query int false "Jumlah item per halaman" default(10)
// @Success      200 {object} tables.PaginatedBookResponse
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/books [GET]
func (c *BookController) GetPublishedBookList(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	books, err := c.bookDAO.GetPublishedBooks(ctx.Context(), limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve book list."})
	}
	totalItems, err := c.bookDAO.CountPublishedBooks(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to count books."})
	}
	totalPages := (totalItems + int64(limit) - 1) / int64(limit)
	response := tables.PaginatedBookResponse{
		Pagination: tables.PaginationInfo{CurrentPage: page, PageSize: limit, TotalItems: totalItems, TotalPages: int(totalPages)},
		Books:      books,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetMyBookDetail adalah handler untuk mendapatkan detail lengkap buku milik penulis yang login.
// @Summary      Dapatkan Detail Buku Saya (Pribadi)
// @Description  Mengambil detail lengkap sebuah buku, termasuk daftar chapternya. Hanya bisa diakses oleh penulis buku tersebut.
// @Tags         Book Management
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookId path int true "ID Buku"
// @Success      200 {object} tables.BookDetailResponse
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      403 {object} ErrorResponse "Bukan pemilik buku"
// @Failure      404 {object} ErrorResponse "Buku tidak ditemukan"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/books/{bookId}/detail [GET]
// GetMyBookDetail adalah handler untuk mendapatkan detail lengkap buku milik penulis yang login.
func (c *BookController) GetMyBookDetail(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid user token."})
	}
	bookId, err := strconv.ParseInt(ctx.Params("bookId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid book ID."})
	}
	book, err := c.bookDAO.GetBookDetailByID(ctx.Context(), bookId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get book details."})
	}
	if book == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: constants.ErrCodeBookNotFound, Message: "Book not found."})
	}
	if book.AuthorID != userId {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Code: constants.ErrCodeBookNotOwner, Message: "You are not the owner of this book."})
	}
	
	// Panggil dengan isPublic = false agar chapter draft juga muncul untuk penulis
	chapters, err := c.chapterDAO.GetChaptersByBookID(ctx.Context(), bookId, false)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get chapters."})
	}
    if chapters == nil {
        chapters = []tables.Chapter{}
    }

	author, err := c.userDAO.FindUserByID(ctx.Context(), userId)
	if err != nil || author == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get author data."})
	}
    author.Password = ""
	
	reviews, err := c.reviewDAO.GetReviewsByBookID(ctx.Context(), bookId)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get reviews."})
    }
    if reviews == nil {
        reviews = []tables.Review{}
    }

	response := tables.BookDetailResponse{
		BookInfo: book,
		Chapters: chapters,
		Author:   author,
		Reviews:  reviews,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// processBookStatus adalah fungsi helper internal untuk validasi umum sebelum mengubah status buku.
func (c *BookController) processBookStatus(ctx *fiber.Ctx) (int64, int64, *tables.Book, error) {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return 0, 0, nil, ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid user token."})
	}
	bookId, err := strconv.ParseInt(ctx.Params("bookId"), 10, 64)
	if err != nil {
		return 0, 0, nil, ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid book ID."})
	}
	book, err := c.bookDAO.GetBookWithAuthor(ctx.Context(), bookId)
	if err != nil {
		return 0, 0, nil, ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get book details."})
	}
	if book == nil {
		return 0, 0, nil, ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: constants.ErrCodeBookNotFound, Message: "Book not found."})
	}
	if book.AuthorID != userId {
		return 0, 0, nil, ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Code: constants.ErrCodeBookNotOwner, Message: "You are not the owner of this book."})
	}
	return userId, bookId, book, nil
}

// PublishBook mempublikasikan sebuah buku.
// @Summary      Publikasikan Buku
// @Description  Mengubah status buku menjadi 'Published'. Memerlukan minimal 1 chapter.
// @Tags         Book Management
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookId path int true "ID Buku"
// @Success      200 {object} object{code=string,message=string}
// @Failure      400 {object} ErrorResponse "Buku tidak memiliki chapter"
// @Router       /v1/books/{bookId}/publish [PATCH]
func (c *BookController) PublishBook(ctx *fiber.Ctx) error {
	_, bookId, _, err := c.processBookStatus(ctx)
	if err != nil {
		return err
	}
	chapterCount, err := c.bookDAO.CountChaptersByBookID(ctx.Context(), bookId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to count chapters."})
	}
	if chapterCount == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBookNoChapters, Message: "Cannot publish a book with no chapters."})
	}
	if err := c.bookDAO.UpdateBookStatus(ctx.Context(), bookId, "P"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeUserUpdateFailed, Message: "Failed to publish book."})
	}
	return ctx.JSON(fiber.Map{"code": "book.publish.success", "message": "Book published successfully."})
}

// UnpublishBook mengembalikan buku ke status draft.
// @Summary      Batalkan Publikasi Buku
// @Description  Mengubah status buku kembali menjadi 'Draft'. Hanya bisa dilakukan pada buku yang sedang 'Published'.
// @Tags         Book Management
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookId path int true "ID Buku"
// @Success      200 {object} object{code=string,message=string}
// @Router       /v1/books/{bookId}/unpublish [PATCH]
func (c *BookController) UnpublishBook(ctx *fiber.Ctx) error {
	_, bookId, book, err := c.processBookStatus(ctx)
	if err != nil {
		return err
	}
	if book.Status != "P" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBookNotPublished, Message: "Only published books can be unpublished."})
	}
	if err := c.bookDAO.UpdateBookStatus(ctx.Context(), bookId, "D"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeUserUpdateFailed, Message: "Failed to unpublish book."})
	}
	return ctx.JSON(fiber.Map{"code": "book.unpublish.success", "message": "Book unpublished successfully."})
}

// CompleteBook menandai buku sebagai selesai.
// @Summary      Selesaikan Buku
// @Description  Mengubah status buku menjadi 'Completed'. Hanya bisa dilakukan pada buku yang sedang 'Published'.
// @Tags         Book Management
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookId path int true "ID Buku"
// @Success      200 {object} object{code=string,message=string}
// @Router       /v1/books/{bookId}/complete [PATCH]
func (c *BookController) CompleteBook(ctx *fiber.Ctx) error {
	_, bookId, book, err := c.processBookStatus(ctx)
	if err != nil {
		return err
	}
	if book.Status != "P" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBookNotPublished, Message: "Only published books can be marked as completed."})
	}
	if err := c.bookDAO.UpdateBookStatus(ctx.Context(), bookId, "C"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeUserUpdateFailed, Message: "Failed to complete book."})
	}
	return ctx.JSON(fiber.Map{"code": "book.complete.success", "message": "Book marked as completed."})
}

// HoldBook menandai buku sebagai ditunda.
// @Summary      Tunda Buku
// @Description  Mengubah status buku menjadi 'On Hold'.
// @Tags         Book Management
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookId path int true "ID Buku"
// @Success      200 {object} object{code=string,message=string}
// @Router       /v1/books/{bookId}/hold [PATCH]
func (c *BookController) HoldBook(ctx *fiber.Ctx) error {
	_, bookId, _, err := c.processBookStatus(ctx)
	if err != nil {
		return err
	}
	if err := c.bookDAO.UpdateBookStatus(ctx.Context(), bookId, "H"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeUserUpdateFailed, Message: "Failed to put book on hold."})
	}
	return ctx.JSON(fiber.Map{"code": "book.hold.success", "message": "Book put on hold successfully."})
}

// GetPublicBookDetail adalah handler untuk mendapatkan detail buku yang bisa diakses siapa saja.
// @Summary      Dapatkan Detail Buku (Publik)
// @Description  Mengambil detail lengkap sebuah buku, termasuk daftar chapter, penulis, dan ulasan.
// @Tags         Book
// @Produce      json
// @Param        bookId path int true "ID Buku"
// @Success      200 {object} tables.BookDetailResponse
// @Failure      404 {object} ErrorResponse "Buku tidak ditemukan"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/books/{bookId} [GET]
func (c *BookController) GetPublicBookDetail(ctx *fiber.Ctx) error {
	bookId, err := strconv.ParseInt(ctx.Params("bookId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Invalid book ID."})
	}

	// 1. Ambil data buku utama
	book, err := c.bookDAO.GetBookDetailByID(ctx.Context(), bookId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get book details."})
	}
	if book == nil || book.Status == "D" { // Sembunyikan jika buku tidak ditemukan atau masih draft
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: constants.ErrCodeBookNotFound, Message: "Book not found or not published."})
	}

	// 2. Ambil daftar chapter yang sudah publish
	chapters, err := c.chapterDAO.GetChaptersByBookID(ctx.Context(), bookId, true)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get chapters."})
	}
	if chapters == nil {
		chapters = []tables.Chapter{}
	}

	// 3. Ambil data penulis
	author, err := c.userDAO.FindUserByID(ctx.Context(), book.AuthorID)
	if err != nil || author == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get author data."})
	}
	author.Password = ""

	// 4. Ambil ulasan untuk buku ini
	reviews, err := c.reviewDAO.GetReviewsByBookID(ctx.Context(), bookId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to get reviews."})
	}
	if reviews == nil {
		reviews = []tables.Review{}
	}

	// 5. Gabungkan semua data menjadi satu response
	response := tables.BookDetailResponse{
		BookInfo: book,
		Chapters: chapters,
		Author:   author,
		Reviews:  reviews,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
