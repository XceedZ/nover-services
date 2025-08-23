package controllers

import (
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	walletDAO *dao.WalletDao
	log       *logrus.Logger
}

func NewWalletController(walletDAO *dao.WalletDao) *WalletController {
	return &WalletController{
		walletDAO: walletDAO,
		log:       logrus.New(),
	}
}

// GetMyWallet adalah handler untuk mengambil saldo koin pengguna yang login.
// @Summary      Dapatkan Saldo Koin Saya
// @Description  Mengambil saldo Koin Biasa dan Koin Bonus milik pengguna yang sedang login.
// @Tags         Wallet
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} tables.Wallet
// @Failure      401 {object} ErrorResponse "Tidak terotentikasi"
// @Failure      500 {object} ErrorResponse "Error internal server"
// @Router       /v1/wallet/my-balance [GET]
func (c *WalletController) GetMyWallet(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok || userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Code: constants.ErrCodeUserUnauthorized, Message: "Invalid access."})
	}

	wallet, err := c.walletDAO.GetWalletByUserID(ctx.Context(), userId)
	if err != nil {
		c.log.WithError(err).Errorf("Gagal mengambil wallet untuk user ID: %d", userId)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve wallet balance."})
	}
	
	return ctx.Status(fiber.StatusOK).JSON(wallet)
}