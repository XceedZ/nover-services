package tables

import "time"

// Book merepresentasikan data dari tabel 'books'.
type Book struct {
	BookID        int64      `json:"bookId" db:"book_id"`
	Title         string     `json:"title" db:"title"`
	Description   *string    `json:"description,omitempty" db:"description"`
	CoverImageURL *string    `json:"coverImageUrl,omitempty" db:"cover_image_url"`
	Status        string     `json:"status" db:"status"`
	RatingAverage float64    `json:"ratingAverage" db:"rating_average"`
	TotalViews    int64      `json:"totalViews" db:"total_views"`
	CreateDatetime time.Time `json:"createDatetime" db:"create_datetime"`
	UpdateDatetime *time.Time `json:"updateDatetime,omitempty" db:"update_datetime"`
    // --- PERUBAHAN DI SINI ---
    // Field baru untuk menampung genre sebagai satu string.
    // 'db:"genres"' akan memetakan hasil dari STRING_AGG.
    Genres        *string    `json:"genres,omitempty" db:"genres"` 
}

// Chapter merepresentasikan data dari tabel 'chapters'.
type Chapter struct {
	ChapterID     int64      `json:"chapterId"`
	BookID        int64      `json:"bookId"`
	Title         string     `json:"title"`
	Content       *string    `json:"content,omitempty"`
	ChapterOrder  int        `json:"chapterOrder"`
	Status        string     `json:"status"`
	CoinCost      int        `json:"coinCost"`
	TotalViews    int64      `json:"totalViews"`
	CreateDatetime time.Time `json:"createDatetime"`
	UpdateDatetime *time.Time `json:"updateDatetime,omitempty"`
}
