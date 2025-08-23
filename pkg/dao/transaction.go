package dao

import (
	"context"
	"fmt"
	"noversystem/pkg/tables"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type TransactionDao struct {
	DB *pgxpool.Pool
}

func NewTransactionDao(db *pgxpool.Pool) *TransactionDao {
	return &TransactionDao{DB: db}
}

// GetTransactionsByUserID mengambil riwayat transaksi berdasarkan tipe (pendapatan/pengeluaran).
func (d *TransactionDao) GetTransactionsByUserID(ctx context.Context, userID int64, isDebit bool, limit, offset int) ([]tables.CoinTransaction, error) {
	var transactions []tables.CoinTransaction

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	conditions := squirrel.And{
		squirrel.Eq{"user_id": userID},
	}
	if isDebit {
		conditions = append(conditions, squirrel.Lt{"amount": 0})
	} else {
		conditions = append(conditions, squirrel.Gt{"amount": 0})
	}

	// ✨ PERBAIKAN: Ganti SELECT * dengan kolom yang spesifik
	queryBuilder := psql.Select(
		"transaction_id",
		"user_id",
		"transaction_type",
		"coin_type",
		"amount",
		"description",
		"expiry_date",
		"create_datetime",
	).From("coin_transactions").
		Where(conditions).
		OrderBy("create_datetime DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("gagal membangun query transaksi: %w", err)
	}

	err = pgxscan.Select(ctx, d.DB, &transactions, sql, args...)
	if err != nil {
		// Tambahkan log ini untuk melihat error database yang lebih spesifik
		logrus.WithError(err).Error("Error saat scanning data transaksi")
		return nil, fmt.Errorf("gagal mengambil transaksi: %w", err)
	}
	return transactions, nil
}
