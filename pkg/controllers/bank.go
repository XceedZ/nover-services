package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type BankController struct {
	DB  *pgxpool.Pool
	Log *logrus.Logger
}

func NewBankController(db *pgxpool.Pool) *BankController {
	return &BankController{DB: db}
}

type BankResponse struct {
	BankId   int64  `json:"bank_id"`
	BankName string `json:"bank_name"`
	BankCode string `json:"bank_code"`
}

// GetBankList adalah handler untuk mengambil daftar bank.
// @Summary      Daftar Bank
// @Description  Mengambil daftar semua bank yang aktif, diurutkan berdasarkan nama.
// @Tags         Bank
// @Accept       json
// @Produce      json
// @Success      200 {array}  BankResponse "Daftar bank yang berhasil diambil"
// @Failure      500 {object} fiber.Map    "Terjadi kesalahan internal pada server"
// @Router       /v1/bank/get [get]
func (c *BankController) GetBankList(ctx *fiber.Ctx) error {
	const getBanksQuery = `
		SELECT 
			bank_id, 
			bank_name, 
			COALESCE(bank_code, '') AS bank_code
		FROM banks
		WHERE non_active_datetime IS NULL
		ORDER BY bank_name ASC;`

	rows, err := c.DB.Query(context.Background(), getBanksQuery)
	if err != nil {
		c.Log.WithError(err).Error("Gagal mengeksekusi query daftar bank")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Terjadi kesalahan pada server",
		})
	}
	defer rows.Close()

	var banks []BankResponse
	for rows.Next() {
		var bank BankResponse
		if err := rows.Scan(&bank.BankId, &bank.BankName, &bank.BankCode); err != nil {
			c.Log.WithError(err).Error("Gagal memindai data bank dari baris query")
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Terjadi kesalahan pada server",
			})
		}
		banks = append(banks, bank)
	}

	if err := rows.Err(); err != nil {
		c.Log.WithError(err).Error("Error saat iterasi hasil query bank")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Terjadi kesalahan pada server",
		})
	}

	if banks == nil {
		banks = make([]BankResponse, 0)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"bankList": banks,
	})
}
