package controllers

import (
	"noversystem/pkg/constants"
	"noversystem/pkg/dao"
	"noversystem/pkg/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	txDAO *dao.TransactionDao
	log   *logrus.Logger
}

func NewTransactionController(txDAO *dao.TransactionDao) *TransactionController {
	return &TransactionController{txDAO: txDAO, log: logrus.New()}
}

// GetMyTransactions mengambil riwayat transaksi pengguna.
// @Summary      Dapatkan Riwayat Transaksi
// @Tags         Wallet
// @Produce      json
// @Security     ApiKeyAuth
// @Param        type query string true "Tipe transaksi ('earn' atau 'spend')"
// @Success      200 {object} []tables.CoinTransaction
// @Router       /v1/wallet/transactions [GET]
func (c *TransactionController) GetMyTransactions(ctx *fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(int64)
	txType := ctx.Query("type", "earn") // default 'earn'
	isDebit := txType == "spend"

	// Di sini Anda bisa menambahkan pagination jika diperlukan
	transactions, err := c.txDAO.GetTransactionsByUserID(ctx.Context(), userId, isDebit, 50, 0) // Limit 50
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Code: constants.ErrCodeInternalServer, Message: "Failed to retrieve transactions."})
	}
	if transactions == nil {
		transactions = make([]tables.CoinTransaction, 0)
	}
	return ctx.Status(fiber.StatusOK).JSON(transactions)
}
