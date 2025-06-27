package dao

import (
	"context"
	"errors"
	"noversystem/pkg/tables"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan" // Pastikan ini v2 untuk pgx/v5
	"github.com/jackc/pgx/v5"                // UBAH KE v5
	"github.com/jackc/pgx/v5/pgxpool"        // UBAH KE v5
)

type UserDao struct {
	DB *pgxpool.Pool // Sekarang ini adalah *pgxpool.Pool dari v5
}

func NewUserDao(db *pgxpool.Pool) *UserDao { // Fungsi ini sekarang menerima *pgxpool.Pool dari v5
	return &UserDao{DB: db}
}

// RegisterUser menyimpan pengguna baru ke database menggunakan Squirrel Query Builder.
func (d *UserDao) RegisterUser(ctx context.Context, user *tables.User) (int64, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Insert("users").
		Columns("user_code", "email", "password", "full_name", "username", "login_with").
		Values(user.UserCode, user.Email, user.Password, user.FullName, user.Username, user.LoginWith).
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

// FindUserByEmail mencari pengguna berdasarkan alamat email.
func (d *UserDao) FindUserByEmail(ctx context.Context, email string) (*tables.User, error) {
	var user tables.User
	
	sql := "SELECT user_id, user_code, email, password, full_name, login_with FROM users WHERE email = $1 AND login_with = 'local'"
	
	err := pgxscan.Get(ctx, d.DB, &user, sql, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindUserByUsername mencari pengguna berdasarkan username.
func (d *UserDao) FindUserByUsername(ctx context.Context, username string) (*tables.User, error) {
    var user tables.User

    // Query diubah untuk mencari berdasarkan username
	sql := `SELECT 
				user_id, user_code, email, password, full_name, username, 
				avatar_url, login_with, is_email_verified, create_datetime, update_datetime
			FROM users 
			WHERE username = $1 AND login_with = 'local'`

    err := pgxscan.Get(ctx, d.DB, &user, sql, username)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            // Pengguna tidak ditemukan, bukan error fatal
            return nil, nil 
        }
        return nil, err // Error database lainnya
    }

    return &user, nil
}