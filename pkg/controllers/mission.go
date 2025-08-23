package controllers

import (
	"noversystem/pkg/dao"
	"github.com/gofiber/fiber/v2"
)

type MissionController struct {
	missionDAO *dao.MissionDao
}

func NewMissionController(missionDAO *dao.MissionDao) *MissionController {
	return &MissionController{missionDAO: missionDAO}
}

// GetDailyMissions mengambil daftar misi harian.
// @Summary      Dapatkan Misi Harian
// @Tags         Event
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} []tables.MissionStatus
// @Router       /v1/missions/daily [GET]
func (c *MissionController) GetDailyMissions(ctx *fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(int64)
	missions, err := c.missionDAO.GetActiveMissionsWithUserProgress(ctx.Context(), userId)
	if err != nil {
		// handle error
		return ctx.SendStatus(500)
	}
	return ctx.JSON(missions)
}