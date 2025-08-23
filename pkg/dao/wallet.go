package dao

import (
	"context"
	"fmt"
	"noversystem/pkg/tables"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletDao struct {
	DB *pgxpool.Pool
}

func NewWalletDao(db *pgxpool.Pool) *WalletDao {
	return &WalletDao{DB: db}
}

// GetWalletByUserID mengambil data wallet untuk seorang pengguna.
// Jika wallet belum ada, fungsi ini akan membuatnya.
func (d *WalletDao) GetWalletByUserID(ctx context.Context, userID int64) (*tables.Wallet, error) {
	var wallet tables.Wallet
	
	// Coba ambil wallet yang ada
	query := `SELECT wallet_id, user_id, paid_coins, bonus_coins, update_datetime FROM wallets WHERE user_id = $1`
	err := pgxscan.Get(ctx, d.DB, &wallet, query, userID)
	
	if err != nil {
		// Jika tidak ada baris (wallet belum ada), buat baru.
		if pgxscan.NotFound(err) {
			insertQuery := `
				INSERT INTO wallets (user_id, paid_coins, bonus_coins) 
				VALUES ($1, 0, 0) 
				RETURNING wallet_id, user_id, paid_coins, bonus_coins, update_datetime
			`
			err = pgxscan.Get(ctx, d.DB, &wallet, insertQuery, userID)
			if err != nil {
				return nil, fmt.Errorf("gagal membuat wallet baru: %w", err)
			}
			return &wallet, nil
		}
		// Error lain
		return nil, fmt.Errorf("gagal mengambil wallet: %w", err)
	}

	return &wallet, nil
}