package dao

import (
	"context"
	"errors"
	"noversystem/pkg/tables"

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

func (d *UserDao) FindUserByUsername(ctx context.Context, username string) (*tables.User, error) {
    var user tables.User

	sql := `SELECT 
				user_id, user_code, email, password, full_name, username, 
				avatar_url, login_with, is_email_verified, create_datetime, update_datetime
			FROM users 
			WHERE username = $1 AND login_with = 'local'`

    err := pgxscan.Get(ctx, d.DB, &user, sql, username)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil 
        }
        return nil, err
    }

    return &user, nil
}