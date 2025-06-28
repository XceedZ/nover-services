package dao

import (
	"context"
	"noversystem/pkg/tables"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GenreDao menangani semua operasi database yang terkait dengan genre.
type GenreDao struct {
	DB *pgxpool.Pool
}

func NewGenreDao(db *pgxpool.Pool) *GenreDao {
	return &GenreDao{DB: db}
}

// GetAllActiveGenres mengambil semua genre yang aktif dari database.
// Genre dianggap aktif jika non_active_datetime adalah NULL.
func (d *GenreDao) GetAllActiveGenres(ctx context.Context) ([]tables.Genre, error) {
	var genres []tables.Genre
	// Asumsi Anda punya kolom display_order untuk pengurutan kustom
	const query = `
        SELECT 
            genre_id, genre_name, genre_tl, remark, active_datetime, 
            non_active_datetime, create_datetime, update_datetime
        FROM genres
        WHERE non_active_datetime IS NULL
        ORDER BY genre_name ASC`

	err := pgxscan.Select(ctx, d.DB, &genres, query)
	if err != nil {
		return nil, err
	}

	return genres, nil
}
