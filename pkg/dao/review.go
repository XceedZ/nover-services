package dao

import (
	"context"
	"fmt"
	"noversystem/pkg/tables" // Pastikan nama modul Go Anda benar

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
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
// (Fungsi ini tidak berubah)
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

// --- FUNGSI BARU DITAMBAHKAN ---

// CreateReviewAndNotify membuat review baru dan mengirim notifikasi ke penulis buku dalam satu transaksi.
func (d *ReviewDao) CreateReviewAndNotify(ctx context.Context, reviewData *tables.Review, actorName, bookTitle string) (*tables.Review, error) {
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Langkah 1: Dapatkan ID penulis buku
	var authorID int64
	authorQuery, authorArgs, _ := psql.Select("user_id").
		From("author_books").
		Where(squirrel.Eq{"book_id": reviewData.BookID}).
		Limit(1).
		ToSql()

	err = tx.QueryRow(ctx, authorQuery, authorArgs...).Scan(&authorID)
	if err != nil {
		logrus.Errorf("Gagal mendapatkan author_id untuk book_id %d: %v", reviewData.BookID, err)
		return nil, fmt.Errorf("buku tidak ditemukan atau tidak memiliki penulis: %w", err)
	}

	// Langkah 2: Sisipkan review baru
	reviewQuery, reviewArgs, _ := psql.Insert("reviews").
		Columns("book_id", "user_id", "rating", "review_text").
		Values(reviewData.BookID, reviewData.UserID, reviewData.Rating, reviewData.ReviewText).
		Suffix("RETURNING review_id, create_datetime").
		ToSql()

	err = tx.QueryRow(ctx, reviewQuery, reviewArgs...).Scan(&reviewData.ReviewID, &reviewData.CreateDatetime)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, fmt.Errorf("anda sudah memberikan review untuk buku ini")
		}
		logrus.Errorf("Gagal INSERT ke tabel reviews: %v", err)
		return nil, fmt.Errorf("gagal menyimpan review: %w", err)
	}

	// Langkah 3: Sisipkan notifikasi (jika penulis bukan orang yang sama dengan yang mereview)
	if authorID != reviewData.UserID {
		notificationContent := fmt.Sprintf("%s memberikan review baru untuk buku Anda '%s'.", actorName, bookTitle)
		notifQuery, notifArgs, _ := psql.Insert("system_notifications").
			Columns("user_id", "actor_id", "notification_type", "content", "related_entity_type", "related_entity_id").
			Values(authorID, reviewData.UserID, "NEW_RATING", notificationContent, "BOOK", reviewData.BookID).
			ToSql()

		if _, err := tx.Exec(ctx, notifQuery, notifArgs...); err != nil {
			logrus.Errorf("Gagal INSERT ke tabel system_notifications: %v", err)
			return nil, fmt.Errorf("gagal membuat notifikasi: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return reviewData, nil
}