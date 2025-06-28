package dao

import (
	"context"
	"noversystem/pkg/tables"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// BookDao menangani semua operasi database yang terkait dengan buku.
type BookDao struct {
	DB *pgxpool.Pool
}

func NewBookDao(db *pgxpool.Pool) *BookDao {
	return &BookDao{DB: db}
}

// CreateBook membuat entri buku baru beserta relasinya dalam satu transaksi.
// (Fungsi ini tidak berubah)
func (d *BookDao) CreateBook(ctx context.Context, bookData *tables.Book, authorID int64, genreIDs []int64) (*tables.Book, error) {
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	var newBookID int64
	sql, args, err := psql.Insert("books").
		Columns("title", "description", "cover_image_url").
		Values(bookData.Title, bookData.Description, bookData.CoverImageURL).
		Suffix("RETURNING book_id").
		ToSql()
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(ctx, sql, args...).Scan(&newBookID)
	if err != nil {
		logrus.Errorf("Gagal INSERT ke tabel books: %v", err)
		return nil, err
	}
	bookData.BookID = newBookID

	sql, args, err = psql.Insert("author_books").
		Columns("user_id", "book_id").
		Values(authorID, newBookID).
		ToSql()
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		logrus.Errorf("Gagal INSERT ke tabel author_books: %v", err)
		return nil, err
	}

	if len(genreIDs) > 0 {
		insertBuilder := psql.Insert("book_genres").Columns("book_id", "genre_id")
		for _, genreID := range genreIDs {
			insertBuilder = insertBuilder.Values(newBookID, genreID)
		}
		sql, args, err = insertBuilder.ToSql()
		if err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			logrus.Errorf("Gagal INSERT ke tabel book_genres: %v", err)
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return bookData, nil
}


// --- FUNGSI GetBooksByAuthorID DIPERBARUI TOTAL ---

// GetBooksByAuthorID mengambil semua buku yang ditulis oleh seorang penulis,
// lengkap dengan daftar genre yang digabungkan menjadi satu string.
func (d *BookDao) GetBooksByAuthorID(ctx context.Context, authorID int64) ([]tables.Book, error) {
	var books []tables.Book

	// Query ini menggunakan LEFT JOIN dan STRING_AGG untuk mendapatkan nama genre
	// dalam satu string, dipisahkan oleh koma dan spasi.
	const query = `
		SELECT
			b.book_id, b.title, b.description, b.cover_image_url, b.status,
			b.rating_average, b.total_views, b.create_datetime, b.update_datetime,
			STRING_AGG(g.genre_name, ', ') as genres
		FROM
			books b
		JOIN
			author_books ab ON b.book_id = ab.book_id
		LEFT JOIN
			book_genres bg ON b.book_id = bg.book_id
		LEFT JOIN
			genres g ON bg.genre_id = g.genre_id
		WHERE
			ab.user_id = $1
		GROUP BY
			b.book_id
		ORDER BY
			b.update_datetime DESC NULLS LAST, b.create_datetime DESC
	`

	err := pgxscan.Select(ctx, d.DB, &books, query, authorID)
	if err != nil {
		return nil, err
	}

	return books, nil
}
