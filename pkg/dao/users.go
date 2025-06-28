package dao

import (
	"context"
	"database/sql" // PENTING: Import untuk menggunakan sql.NullString
	"errors"
	"noversystem/pkg/tables" // Pastikan path ini sesuai dengan struktur proyek Anda

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDao struct {
	DB *pgxpool.Pool
}

func NewUserDao(db *pgxpool.Pool) *UserDao {
	return &UserDao{DB: db}
}

// DIPERBAIKI: Menggunakan squirrel dan sql.NullString untuk menangani kolom UNIQUE yang opsional
func (d *UserDao) RegisterUser(ctx context.Context, user *tables.User) (int64, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Menggunakan sql.NullString untuk kolom opsional yang UNIQUE (seperti pen_name dan phone)
	// agar bisa di-INSERT sebagai NULL, bukan string kosong.
	var penName, phone, instagram sql.NullString
	if user.PenName != nil {
		penName = sql.NullString{String: *user.PenName, Valid: true}
	}
	if user.Phone != nil {
		phone = sql.NullString{String: *user.Phone, Valid: true}
	}
	if user.Instagram != nil {
		instagram = sql.NullString{String: *user.Instagram, Valid: true}
	}

	sql, args, err := psql.Insert("users").
		Columns(
			"user_code", "email", "password", "full_name", "username", "login_with",
			"avatar_url", "is_email_verified", "flg_author",
			// Kolom opsional yang unik
			"pen_name", "phone", "instagram",
		).
		Values(
			user.UserCode, user.Email, user.Password, user.FullName, user.Username, user.LoginWith,
			"",    // avatar_url default
			false, // is_email_verified default
			"N",   // flg_author default
			// Memberikan nilai NULL jika tidak ada input
			penName,
			phone,
			instagram,
		).
		Suffix("RETURNING user_id").
		ToSql()

	if err != nil {
		return 0, err
	}

	var newUserID int64
	err = d.DB.QueryRow(ctx, sql, args...).Scan(&newUserID)
	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// DIPERBAIKI: Menggunakan SELECT * untuk memastikan semua data pengguna terambil
func (d *UserDao) FindUserByEmail(ctx context.Context, email string) (*tables.User, error) {
	var user tables.User
	const sql = "SELECT * FROM users WHERE email = $1 AND login_with = 'local'"
	err := pgxscan.Get(ctx, d.DB, &user, sql, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// DIPERBAIKI: Menggunakan SELECT * agar konsisten dan lengkap
func (d *UserDao) FindUserByUsername(ctx context.Context, username string) (*tables.User, error) {
	var user tables.User
	const sql = `SELECT * FROM users WHERE username = $1 AND login_with = 'local'`
	err := pgxscan.Get(ctx, d.DB, &user, sql, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

type AuthorUpdateRequest struct {
	PenName       string
	Phone         string
	Instagram     string
	BankId        int64
	AccountNumber string
}

func (d *UserDao) IsPenNameTaken(ctx context.Context, penName string) (bool, error) {
	const query = `SELECT 1 FROM users WHERE pen_name = $1;`
	var exists int
	err := d.DB.QueryRow(ctx, query, penName).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *UserDao) UpdateUserToAuthor(ctx context.Context, userId int64, params AuthorUpdateRequest) error {
	const query = `
		UPDATE users SET
			pen_name = $1,
			phone = $2,
			instagram = $3,
			bank_id = $4,
			account_number = $5,
			flg_author = 'Y',
			update_datetime = NOW()
		WHERE user_id = $6;
	`
	cmdTag, err := d.DB.Exec(ctx, query,
		params.PenName,
		params.Phone,
		params.Instagram,
		params.BankId,
		params.AccountNumber,
		userId,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != 1 {
		return errors.New("pengguna tidak ditemukan atau tidak ada data yang diperbarui")
	}
	return nil
}

func (d *UserDao) FindUserByID(ctx context.Context, userID int64) (*tables.User, error) {
	var user tables.User
	sql := `
        SELECT 
            user_id, user_code, email, password, full_name, username, pen_name,
            avatar_url, login_with, is_email_verified, phone, instagram,
            bank_id, account_number, flg_author, create_datetime, update_datetime
        FROM users 
        WHERE user_id = $1`

	err := pgxscan.Get(ctx, d.DB, &user, sql, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Pengguna dengan ID tersebut tidak ditemukan.
			return nil, nil
		}
		// Error database lainnya
		return nil, err
	}

	return &user, nil
}