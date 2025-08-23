// File: pkg/controllers/checkin_controller.go
package controllers

import (
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"
	"noversystem/pkg/tables"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CheckinController struct {
	checkinDAO *dao.CheckinDao
	log        *logrus.Logger
}

func NewCheckinController(checkinDAO *dao.CheckinDao) *CheckinController {
	return &CheckinController{checkinDAO: checkinDAO, log: logrus.New()}
}

// GetStatus mengambil status check-in terkini.
// @Summary      Dapatkan Status Check-in
// @Tags         Event
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} tables.CheckinStatusResponse
// @Router       /v1/check-in/status [GET]
func (c *CheckinController) GetStatus(ctx *fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(int64)

	now := time.Now()

	rewards, err := c.checkinDAO.GetAllCheckinRewards(ctx.Context())
	if err != nil {
		c.log.WithError(err).Error("Gagal mengambil konfigurasi hadiah checkin")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Gagal memuat data event."})
	}

	checkedInDates, err := c.checkinDAO.GetUserCheckinsForMonth(ctx.Context(), userId, now.Year(), now.Month())
	if err != nil {
		c.log.WithError(err).Error("Gagal mengambil riwayat checkin")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Gagal memuat data event."})
	}

	todayStr := now.Format("2006-01-02")
	todayCheckedIn := false
	for _, dateStr := range checkedInDates {
		if dateStr == todayStr {
			todayCheckedIn = true
			break
		}
	}

	// Total check-in bulan ini menentukan progres hari ke berapa
	totalCheckinsThisMonth := len(checkedInDates)

	response := tables.CheckinStatusResponse{
		// Field 'ConsecutiveStreak' sekarang kita artikan sebagai total checkin bulan ini
		TotalCheckinsThisMonth: totalCheckinsThisMonth,
		TodayCheckedIn:         todayCheckedIn,
		Rewards:                rewards,
		CheckedInDates:         checkedInDates,
	}

	return ctx.JSON(response)
}

// CheckIn melakukan aksi check-in.
// @Summary      Lakukan Check-in Harian
// @Tags         Event
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} object{message=string, reward=int}
// @Router       /v1/check-in [POST]
func (c *CheckinController) CheckIn(ctx *fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(int64)
	reward, err := c.checkinDAO.PerformCheckin(ctx.Context(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "sudah check-in") {
			return ctx.Status(fiber.StatusConflict).JSON(ErrorResponse{Code: "checkin.already_done", Message: "Anda sudah check-in hari ini."})
		}
		c.log.WithError(err).Error("Gagal melakukan check-in")
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Gagal melakukan check-in."})
	}
	return ctx.JSON(fiber.Map{"message": "Check-in berhasil!", "reward": reward})
}
