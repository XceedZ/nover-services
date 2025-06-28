package controllers

import (
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"
	"noversystem/pkg/tables"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// UserController menangani logika yang berhubungan dengan pengguna.
type UserController struct {
	userDAO *dao.UserDao
	log     *logrus.Logger
}

// NewUserController membuat instance baru dari UserController.
func NewUserController(userDAO *dao.UserDao) *UserController {
	return &UserController{
		userDAO: userDAO,
		log:     logrus.New(), // Inisialisasi logger
	}
}

// AuthorRequestPayload adalah struktur untuk data permintaan menjadi penulis.
type AuthorRequestPayload struct {
	PenName       string `json:"penName"`
	Phone         string `json:"phone"`
	Instagram     string `json:"instagram"`
	BankId        int64  `json:"bankId"`
	AccountNumber string `json:"accountNumber"`
}

// RequestBecomeAuthor adalah handler untuk permintaan menjadi penulis.
// @Summary      Permintaan menjadi Penulis
// @Description  Memperbarui profil pengguna untuk menjadi penulis dengan melengkapi data yang diperlukan. Endpoint ini memerlukan otentikasi.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        author_request body AuthorRequestPayload true "Data untuk menjadi penulis"
// @Success      200 {object} object{code=string,message=string} "Pesan sukses"
// @Failure      400 {object} ErrorResponse "Input tidak valid"
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      409 {object} ErrorResponse "Nama pena sudah digunakan"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/user/request-author [POST]
func (c *UserController) RequestBecomeAuthor(ctx *fiber.Ctx) error {
	// 1. Ambil userID yang sudah divalidasi dari middleware
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Code:    constants.ErrCodeUserUnauthorized,
			Message: "Invalid access, user not authenticated properly.",
		})
	}

	// 2. Parse dan sanitasi input dari body
	var payload AuthorRequestPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    constants.ErrCodeBadRequest,
			Message: "Cannot parse request body.",
		})
	}

	payload.PenName = strings.TrimSpace(payload.PenName)
	payload.Phone = strings.TrimSpace(payload.Phone)
	payload.AccountNumber = strings.TrimSpace(payload.AccountNumber)

	// 3. Validasi input yang wajib diisi
	if payload.PenName == "" || payload.Phone == "" || payload.BankId == 0 || payload.AccountNumber == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    constants.ErrCodeAuthInputRequired,
			Message: "Pen name, phone, bank, and account number fields are required.",
		})
	}

	// 4. Cek apakah nama pena sudah digunakan
	isTaken, err := c.userDAO.IsPenNameTaken(ctx.Context(), payload.PenName)
	if err != nil {
		c.log.WithError(err).Error("Failed to check pen name availability")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    constants.ErrCodeInternalServer,
			Message: "A server error occurred while checking pen name.",
		})
	}
	if isTaken {
		return ctx.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Code:    constants.ErrCodeUserPenNameTaken,
			Message: "This pen name is already taken, please choose another one.",
		})
	}

	// 5. Siapkan data untuk update dan panggil DAO
	updateParams := dao.AuthorUpdateRequest{
		PenName:       payload.PenName,
		Phone:         payload.Phone,
		Instagram:     payload.Instagram,
		BankId:        payload.BankId,
		AccountNumber: payload.AccountNumber,
	}

	if err := c.userDAO.UpdateUserToAuthor(ctx.Context(), userId, updateParams); err != nil {
		c.log.WithError(err).Error("Failed to update user to author")
		// Cek apakah error karena pengguna tidak ditemukan
		if err.Error() == "pengguna tidak ditemukan atau tidak ada data yang diperbarui" {
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Code:    constants.ErrCodeUserNotFound,
				Message: err.Error(),
			})
		}
		// Error server lainnya
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    constants.ErrCodeUserUpdateFailed,
			Message: "Failed to update user profile.",
		})
	}

	// 6. Kirim response sukses
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    "author_request_success",
		"message": "Request to become an author was successful, your profile has been updated.",
	})
}

type AuthorStatusResponse struct {
	IsAuthor bool         `json:"isAuthor"`
	User     *tables.User `json:"user"`
}

// CheckAuthorStatus memeriksa status penulis dan mengembalikan data profil lengkap.
// @Summary      Cek Status Penulis & Profil
// @Description  Memvalidasi token, memeriksa apakah pengguna adalah penulis, dan mengembalikan data profil lengkap pengguna.
// @Tags         User
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} AuthorStatusResponse "Status penulis dan data profil lengkap"
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      404 {object} ErrorResponse "Pengguna tidak ditemukan"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/user/author-status [GET]
func (c *UserController) CheckAuthorStatus(ctx *fiber.Ctx) error {
	// 1. Ambil userID yang sudah divalidasi dari middleware
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Code:    constants.ErrCodeUserUnauthorized,
			Message: "Invalid access, user ID not found in token.",
		})
	}

	// 2. Panggil DAO untuk mengambil data pengguna lengkap berdasarkan ID
	user, err := c.userDAO.FindUserByID(ctx.Context(), userId)
	if err != nil {
		c.log.WithError(err).Errorf("Failed to find user by ID %d", userId)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    constants.ErrCodeInternalServer,
			Message: "Failed to retrieve user data.",
		})
	}
	// Jika pengguna tidak ditemukan di database (kemungkinan aneh jika token valid)
	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Code:    constants.ErrCodeUserNotFound,
			Message: "User associated with this token not found.",
		})
	}

	// 3. Hapus hash password sebelum mengirim response
	user.Password = ""

	// 4. Siapkan dan kembalikan response lengkap
	response := AuthorStatusResponse{
		IsAuthor: user.FlgAuthor == "Y",
		User:     user,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
