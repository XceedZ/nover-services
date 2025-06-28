package controllers

import (
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"
	"noversystem/pkg/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// GenreController menangani logika HTTP yang terkait dengan genre.
type GenreController struct {
	genreDAO *dao.GenreDao
	log      *logrus.Logger
}

// NewGenreController membuat instance baru dari GenreController.
func NewGenreController(genreDAO *dao.GenreDao) *GenreController {
	return &GenreController{
		genreDAO: genreDAO,
		log:      logrus.New(),
	}
}

// --- STRUCT BARU UNTUK RESPONSE ---
// GenreListResponse adalah struktur untuk response daftar genre yang dibungkus.
type GenreListResponse struct {
	GenreList []tables.Genre `json:"genreList"`
}

// GetAllGenres adalah handler untuk mendapatkan semua genre yang aktif.
// @Summary      Dapatkan Semua Genre Aktif
// @Description  Mengambil daftar semua genre yang tersedia dan aktif di sistem.
// @Tags         Genre
// @Produce      json
// @Success      200 {object} GenreListResponse
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/genres [GET]
func (c *GenreController) GetAllGenres(ctx *fiber.Ctx) error {
	genres, err := c.genreDAO.GetAllActiveGenres(ctx.Context())
	if err != nil {
		c.log.WithError(err).Error("Gagal mengambil daftar genre dari DAO")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    constants.ErrCodeInternalServer,
			Message: "Failed to retrieve genre list.",
		})
	}

	// Jika tidak ada genre, kembalikan array kosong di dalam objek.
	if genres == nil {
		genres = []tables.Genre{}
	}

	// Bungkus slice 'genres' di dalam struct GenreListResponse
	response := GenreListResponse{GenreList: genres}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
