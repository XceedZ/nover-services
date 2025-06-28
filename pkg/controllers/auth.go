package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"
	"noversystem/pkg/tables"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type LoginSuccessResponse struct {
	Token string       `json:"token"`
	User  *tables.User `json:"user"`
}

type AuthController struct {
	UserDao *dao.UserDao
}

func NewAuthController(userDao *dao.UserDao) *AuthController {
	return &AuthController{UserDao: userDao}
}

type LoginRequest struct {
	Username string `json:"username" example:"testuser"`
	Password string `json:"password" example:"password123"`
}

// Register adalah handler untuk membuat akun pengguna baru.
// @Summary Registrasi pengguna baru
// @Description Membuat akun pengguna baru dengan email, username, dan password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body tables.User true "Informasi Registrasi Pengguna"
// @Success 201 {object} tables.User
// @Failure 400 {object} ErrorResponse "Input tidak valid"
// @Failure 409 {object} ErrorResponse "Email atau Username sudah terdaftar"
// @Failure 500 {object} ErrorResponse "Error internal server"
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req tables.User
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Cannot parse JSON"})
	}

	if req.Email == "" || req.Password == "" || req.Username == nil || *req.Username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeAuthInputRequired, Message: "Email, password, and username are required"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to hash password"})
	}
	req.Password = string(hashedPassword)

	uniqueString := req.Email + time.Now().String()
	hash := md5.Sum([]byte(uniqueString))
	req.UserCode = hex.EncodeToString(hash[:])
	req.LoginWith = "local"

	newID, err := c.UserDao.RegisterUser(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusConflict).JSON(ErrorResponse{Code: constants.ErrCodeAuthEmailOrUsernameTaken, Message: "Failed to register user, email or username may already exist"})
	}
	req.UserId = newID
	req.Password = ""

	return ctx.Status(fiber.StatusCreated).JSON(req)
}

// Login adalah handler yang sudah diubah untuk menggunakan error codes.
// @Summary Login pengguna
// @Description Mengotentikasi pengguna dengan username dan password, lalu memberikan JWT.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Kredensial Login dengan Username"
// @Success 200 {object} LoginSuccessResponse
// @Failure 400 {object} ErrorResponse "Input tidak valid"
// @Failure 401 {object} ErrorResponse "Kredensial tidak valid"
// @Failure 500 {object} ErrorResponse "Error internal server"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: constants.ErrCodeBadRequest, Message: "Cannot parse JSON"})
	}

	user, err := c.UserDao.FindUserByUsername(ctx.Context(), req.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Database error"})
	}
	if user == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeAuthInvalidCredentials, Message: "Invalid username or password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeAuthInvalidCredentials, Message: "Invalid username or password"})
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeAuthJWTSecretMissing, Message: "JWT secret not configured"})
	}

	claims := jwt.MapClaims{
		"userId":    user.UserId,
		"userCode":  user.UserCode,
		"email":     user.Email,
		"loginWith": user.LoginWith,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeAuthTokenCreation, Message: "Failed to create token"})
	}

	user.Password = ""

	response := LoginSuccessResponse{
		Token: signedToken,
		User:  user,
	}

	return ctx.JSON(response)
}
