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

type BookCommentDao struct {
	DB *pgxpool.Pool
}

func NewBookCommentDao(db *pgxpool.Pool) *BookCommentDao {
	return &BookCommentDao{DB: db}
}

// CreateCommentAndNotify membuat komentar buku dan mengirim notifikasi ke penulis.
func (d *BookCommentDao) CreateCommentAndNotify(ctx context.Context, commentData *tables.BookComment, actorName, bookTitle string) (*tables.BookComment, error) {
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// 1. Simpan komentar baru
	commentQuery, args, _ := psql.Insert("book_comments").
		Columns("book_id", "user_id", "comment_text", "parent_comment_id").
		Values(commentData.BookID, commentData.UserID, commentData.CommentText, commentData.ParentCommentID).
		Suffix("RETURNING comment_id, create_datetime").
		ToSql()

	err = tx.QueryRow(ctx, commentQuery, args...).Scan(&commentData.CommentID, &commentData.CreateDatetime)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan komentar: %w", err)
	}

	// 2. Dapatkan ID penulis buku
	var authorID int64
	err = tx.QueryRow(ctx, "SELECT user_id FROM author_books WHERE book_id = $1 LIMIT 1", commentData.BookID).Scan(&authorID)
	if err != nil {
		// Jika penulis tidak ditemukan, kita tetap commit komentarnya tapi tidak kirim notif.
		logrus.Errorf("Gagal menemukan penulis untuk buku %d, notifikasi tidak dikirim: %v", commentData.BookID, err)
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return nil, fmt.Errorf("gagal commit transaksi setelah gagal cari penulis: %w", commitErr)
		}
		return commentData, nil // Kembalikan komentar yang berhasil dibuat
	}

	// 3. Kirim notifikasi jika komentator bukan si penulis
	if authorID != commentData.UserID {
		notificationContent := fmt.Sprintf("%s meninggalkan komentar di buku Anda '%s'.", actorName, bookTitle)

		// Perbaikan: Simpan referensi ke KOMENTAR, bukan ke BUKU
		notifQuery, notifArgs, _ := psql.Insert("system_notifications").
			Columns("user_id", "actor_id", "notification_type", "content", "related_entity_type", "related_entity_id").
			Values(authorID, commentData.UserID, "NEW_COMMENT", notificationContent, "BOOK_COMMENT", commentData.CommentID).
			ToSql()

		if _, err := tx.Exec(ctx, notifQuery, notifArgs...); err != nil {
			// Jika notif gagal, kita tetap anggap berhasil karena komentar sudah masuk.
			logrus.Errorf("Gagal membuat notifikasi untuk komentar buku: %v", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return commentData, nil
}

func (d *BookCommentDao) GetCommentsByBookID(ctx context.Context, bookID int64) ([]tables.BookComment, error) {
	var comments []tables.BookComment

	// ✨ QUERY DIPERBARUI DENGAN LOGIKA KONDISIONAL
	query := `
        SELECT
            bc.comment_id,
            bc.book_id,
            bc.user_id,
            bc.comment_text,
            bc.parent_comment_id,
            bc.create_datetime,
            u.avatar_url,
            -- Logika IF/ELSE di SQL untuk menentukan nama yang ditampilkan
            CASE 
                WHEN bc.user_id = ab.user_id THEN u.pen_name 
                ELSE u.full_name 
            END as pen_name
        FROM
            book_comments bc
        -- Join ke tabel users untuk mendapatkan detail komentator
        JOIN
            users u ON bc.user_id = u.user_id
        -- Join ke tabel author_books untuk mendapatkan ID penulis buku
        LEFT JOIN
            author_books ab ON bc.book_id = ab.book_id
        WHERE
            bc.book_id = $1
        ORDER BY
            bc.create_datetime ASC`

	err := pgxscan.Select(ctx, d.DB, &comments, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil komentar buku: %w", err)
	}
	return comments, nil
}
