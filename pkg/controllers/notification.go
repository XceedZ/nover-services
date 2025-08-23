package controllers

import (
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"
	"noversystem/pkg/tables"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type NotificationController struct {
	notifDAO *dao.NotificationDao
	log      *logrus.Logger
}

func NewNotificationController(notifDAO *dao.NotificationDao) *NotificationController {
	return &NotificationController{
		notifDAO: notifDAO,
		log:      logrus.New(),
	}
}

// GetNotifications adalah handler untuk mengambil daftar notifikasi milik pengguna yang login.
// @Summary      Dapatkan Notifikasi Saya
// @Description  Mengambil daftar notifikasi untuk pengguna yang sedang login dengan pagination.
// @Tags         Notification
// @Produce      json
// @Security     ApiKeyAuth
// @Param        page query int false "Nomor Halaman" default(1)
// @Param        limit query int false "Jumlah item per halaman" default(10)
// @Success      200 {object} tables.PaginatedNotificationResponse
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/notifications [GET]
func (c *NotificationController) GetNotifications(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid access."})
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 { page = 1 }
	if limit > 50 { limit = 50 }
	offset := (page - 1) * limit

	notifications, err := c.notifDAO.GetNotificationsByUserID(ctx.Context(), userId, limit, offset)
	if err != nil {
		c.log.WithError(err).Error("Gagal mengambil notifikasi dari DAO")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve notifications."})
	}
	
	totalItems, err := c.notifDAO.CountNotificationsByUserID(ctx.Context(), userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to count notifications."})
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)
	response := tables.PaginatedNotificationResponse{
		Pagination: tables.PaginationInfo{CurrentPage: page, PageSize: limit, TotalItems: totalItems, TotalPages: int(totalPages)},
		Items:      notifications,
	}

	if response.Items == nil {
		response.Items = make([]tables.NotificationResponse, 0)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}