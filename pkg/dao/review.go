package dao

import (
	"context"
	"noversystem/pkg/tables"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReviewDao menangani operasi database untuk tabel 'reviews'.
type ReviewDao struct {
	DB *pgxpool.Pool
}

// NewReviewDao membuat instance baru dari ReviewDao.
func NewReviewDao(db *pgxpool.Pool) *ReviewDao {
	return &ReviewDao{DB: db}
}

// GetReviewsByBookID mengambil daftar ulasan untuk sebuah buku, lengkap dengan data penulisnya.
func (d *ReviewDao) GetReviewsByBookID(ctx context.Context, bookID int64) ([]tables.Review, error) {
	var reviews []tables.Review
	const query = `
		SELECT
			r.review_id,
			r.user_id,
			u.pen_name,
			u.avatar_url,
			r.rating,
			r.review_text,
			r.create_datetime
		FROM
			reviews r
		JOIN
			users u ON r.user_id = u.user_id
		WHERE
			r.book_id = $1
		ORDER BY
			r.create_datetime DESC`

	err := pgxscan.Select(ctx, d.DB, &reviews, query, bookID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
