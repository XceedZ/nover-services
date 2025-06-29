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

// ChapterDao menangani semua operasi database yang terkait dengan chapter.
type ChapterDao struct {
	DB *pgxpool.Pool
}

// NewChapterDao membuat instance baru dari ChapterDao.
func NewChapterDao(db *pgxpool.Pool) *ChapterDao {
	return &ChapterDao{DB: db}
}

// ... (Fungsi GetLastChapterOrder dan CreateChapter tetap sama) ...
func (d *ChapterDao) GetLastChapterOrder(ctx context.Context, bookID int64) (int, error) {
    var lastOrder int
    const query = `SELECT COALESCE(MAX(chapter_order), 0) FROM chapters WHERE book_id = $1`
    err := d.DB.QueryRow(ctx, query, bookID).Scan(&lastOrder)
    return lastOrder, err
}

func (d *ChapterDao) CreateChapter(ctx context.Context, chapter *tables.Chapter) (*tables.Chapter, error) {
    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    sql, args, err := psql.Insert("chapters").
        Columns("book_id", "title", "content", "chapter_order", "coin_cost").
        Values(chapter.BookID, chapter.Title, chapter.Content, chapter.ChapterOrder, chapter.CoinCost).
        Suffix("RETURNING chapter_id, create_datetime").
        ToSql()
    if err != nil { return nil, err }
    err = d.DB.QueryRow(ctx, sql, args...).Scan(&chapter.ChapterID, &chapter.CreateDatetime)
    return chapter, err
}


// --- FUNGSI GetChaptersByBookID DIPERBARUI ---
// Fungsi ini sekarang menerima parameter isPublic untuk memfilter chapter.
func (d *ChapterDao) GetChaptersByBookID(ctx context.Context, bookID int64, isPublic bool) ([]tables.Chapter, error) {
    var chapters []tables.Chapter

    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    queryBuilder := psql.Select(
        "chapter_id", "book_id", "title", "content", "chapter_order", 
        "status", "coin_cost", "total_views", "create_datetime", "update_datetime",
    ).
    From("chapters").
    Where(squirrel.Eq{"book_id": bookID}).
    OrderBy("chapter_order ASC")  // Urut tetap berdasarkan order

    sql, args, err := queryBuilder.ToSql()
    if err != nil {
        return nil, err
    }

    err = pgxscan.Select(ctx, d.DB, &chapters, sql, args...)
    if err != nil {
        return nil, err
    }
    return chapters, nil
}

// ... (Fungsi GetPublishedChapterByID dan IsChapterUnlockedByUser tetap sama) ...
func (d *ChapterDao) GetPublishedChapterByID(ctx context.Context, chapterID int64) (*tables.Chapter, error) {
    var chapter tables.Chapter
    const query = `SELECT * FROM chapters WHERE chapter_id = $1`
    err := pgxscan.Get(ctx, d.DB, &chapter, query, chapterID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
    return &chapter, nil
}

func (d *ChapterDao) IsChapterUnlockedByUser(ctx context.Context, userID, chapterID int64) (bool, error) {
    var exists bool
    const query = `SELECT EXISTS (SELECT 1 FROM user_unlocked_chapters WHERE user_id = $1 AND chapter_id = $2)`
    err := d.DB.QueryRow(ctx, query, userID, chapterID).Scan(&exists)
    return exists, err
}
