package dao

import (
	"context"
	"errors"
	"noversystem/pkg/tables"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
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

// GetBooksByAuthorID sekarang memiliki parameter untuk membedakan panggilan publik dan pribadi
func (d *BookDao) GetBooksByAuthorID(ctx context.Context, authorID int64, isPublic bool) ([]tables.Book, error) {
	var books []tables.Book

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := psql.Select(
		"b.book_id", "b.title", "b.description", "b.cover_image_url", "b.status",
		"b.rating_average", "b.total_views", "b.create_datetime", "b.update_datetime",
		"STRING_AGG(g.genre_name, ', ') as genres",
	).
		From("books b").
		Join("author_books ab ON b.book_id = ab.book_id").
		LeftJoin("book_genres bg ON b.book_id = bg.book_id").
		LeftJoin("genres g ON bg.genre_id = g.genre_id").
		Where(squirrel.Eq{"ab.user_id": authorID}).
		GroupBy("b.book_id").
		OrderBy("b.update_datetime DESC NULLS LAST", "b.create_datetime DESC")
	
	// Jika panggilan ini untuk publik, tambahkan filter status
	if isPublic {
		queryBuilder = queryBuilder.Where(squirrel.NotEq{"b.status": "D"}) // D = Draft
	}
	
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = pgxscan.Select(ctx, d.DB, &books, sql, args...)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (d *BookDao) GetBookWithAuthor(ctx context.Context, bookID int64) (*tables.Book, error) {
	var book tables.Book
	const query = `
		SELECT
			b.book_id, b.title, b.status, ab.user_id as author_id
		FROM
			books b
		JOIN
			author_books ab ON b.book_id = ab.book_id
		WHERE
			b.book_id = $1
		LIMIT 1;` // Hanya mengambil satu penulis utama untuk validasi

	err := pgxscan.Get(ctx, d.DB, &book, query, bookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Buku tidak ditemukan
		}
		return nil, err
	}
	return &book, nil
}

// CountChaptersByBookID menghitung jumlah chapter yang ada untuk sebuah buku.
func (d *BookDao) CountChaptersByBookID(ctx context.Context, bookID int64) (int, error) {
	var count int
	const query = `SELECT COUNT(*) FROM chapters WHERE book_id = $1`
	err := d.DB.QueryRow(ctx, query, bookID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateBookStatus mengubah status sebuah buku.
func (d *BookDao) UpdateBookStatus(ctx context.Context, bookID int64, newStatus string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := psql.Update("books").
		Set("status", newStatus).
		Where(squirrel.Eq{"book_id": bookID}).
		ToSql()

	if err != nil {
		return err
	}

	cmdTag, err := d.DB.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != 1 {
		return errors.New("buku tidak ditemukan atau status tidak berubah")
	}
	return nil
}

// GetBookDetailByID mengambil detail buku tunggal, lengkap dengan genre yang digabungkan.
func (d *BookDao) GetBookDetailByID(ctx context.Context, bookID int64) (*tables.Book, error) {
    var book tables.Book
	const query = `
		SELECT
			b.book_id, b.title, b.description, b.cover_image_url, b.status,
			b.rating_average, b.total_views, b.create_datetime, b.update_datetime,
			ab.user_id as author_id, -- Sertakan author_id untuk validasi
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
			b.book_id = $1
		GROUP BY
			b.book_id, ab.user_id
	`
	err := pgxscan.Get(ctx, d.DB, &book, query, bookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Buku tidak ditemukan
		}
		return nil, err
	}
	return &book, nil
}

// GetPublishedBooks mengambil daftar buku yang statusnya bukan Draft dengan pagination.
func (d *BookDao) GetPublishedBooks(ctx context.Context, limit, offset int) ([]tables.Book, error) {
    var books []tables.Book
    const query = `
        SELECT
            b.book_id, b.title, b.description, b.cover_image_url, b.status,
            b.rating_average, b.total_views, b.create_datetime, b.update_datetime,
            u.pen_name,
            STRING_AGG(g.genre_name, ', ') as genres
        FROM
            books b
        JOIN
            author_books ab ON b.book_id = ab.book_id
        JOIN
            users u ON ab.user_id = u.user_id
        LEFT JOIN
            book_genres bg ON b.book_id = bg.book_id
        LEFT JOIN
            genres g ON bg.genre_id = g.genre_id
        WHERE
            b.status <> 'D' -- PERUBAHAN: Mengambil semua yang BUKAN Draft ('P', 'C', 'H')
        GROUP BY
            b.book_id, u.pen_name
        ORDER BY
            b.create_datetime DESC
        LIMIT $1 OFFSET $2`

    err := pgxscan.Select(ctx, d.DB, &books, query, limit, offset)
    return books, err
}

// CountPublishedBooks menghitung total buku yang statusnya bukan Draft.
func (d *BookDao) CountPublishedBooks(ctx context.Context) (int64, error) {
    var count int64
    const query = `SELECT COUNT(*) FROM books WHERE status <> 'D'` // PERUBAHAN
    err := d.DB.QueryRow(ctx, query).Scan(&count)
    return count, err
}