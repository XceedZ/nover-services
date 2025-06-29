package dao

import (
	"context"
	"noversystem/pkg/tables"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChapterDao menangani semua operasi database yang terkait dengan chapter.
type ChapterDao struct {
	DB *pgxpool.Pool
}

func NewChapterDao(db *pgxpool.Pool) *ChapterDao {
	return &ChapterDao{DB: db}
}

// GetLastChapterOrder mendapatkan nomor urut chapter terakhir dari sebuah buku.
func (d *ChapterDao) GetLastChapterOrder(ctx context.Context, bookID int64) (int, error) {
	var lastOrder int
	const query = `SELECT COALESCE(MAX(chapter_order), 0) FROM chapters WHERE book_id = $1`

	err := d.DB.QueryRow(ctx, query, bookID).Scan(&lastOrder)
	if err != nil {
		// pgx.ErrNoRows tidak akan terjadi karena MAX/COALESCE selalu mengembalikan satu baris.
		return 0, err
	}
	return lastOrder, nil
}

// CreateChapter menyisipkan sebuah chapter baru ke dalam database.
func (d *ChapterDao) CreateChapter(ctx context.Context, chapter *tables.Chapter) (*tables.Chapter, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Insert("chapters").
		Columns("book_id", "title", "content", "chapter_order", "coin_cost").
		Values(chapter.BookID, chapter.Title, chapter.Content, chapter.ChapterOrder, chapter.CoinCost).
		Suffix("RETURNING chapter_id, create_datetime"). // Ambil kembali data yang digenerate DB
		ToSql()

	if err != nil {
		return nil, err
	}

	// Scan ID dan create_datetime yang baru dibuat ke dalam objek chapter
	err = d.DB.QueryRow(ctx, sql, args...).Scan(&chapter.ChapterID, &chapter.CreateDatetime)
	if err != nil {
		return nil, err
	}

	return chapter, nil
}

// GetChaptersByBookID mengambil semua chapter dari sebuah buku, diurutkan berdasarkan chapter_order.
func (d *ChapterDao) GetChaptersByBookID(ctx context.Context, bookID int64) ([]tables.Chapter, error) {
	var chapters []tables.Chapter
	const query = `
        SELECT 
            chapter_id, book_id, title, content, chapter_order, 
            status, coin_cost, total_views, create_datetime, update_datetime
        FROM chapters
        WHERE book_id = $1
        ORDER BY chapter_order ASC`

	err := pgxscan.Select(ctx, d.DB, &chapters, query, bookID)
	if err != nil {
		return nil, err
	}
	return chapters, nil
}
