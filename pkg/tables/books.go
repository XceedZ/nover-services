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
	CreateDatetime time.Time  `json:"createDatetime" db:"create_datetime"`
	UpdateDatetime *time.Time `json:"updateDatetime,omitempty" db:"update_datetime"`
    Genres        *string    `json:"genres,omitempty" db:"genres"` 
    AuthorID      int64      `json:"-" db:"author_id"`
}

// Chapter merepresentasikan data dari tabel 'chapters'.
type Chapter struct {
	ChapterID     int64      `json:"chapterId" db:"chapter_id"`
	BookID        int64      `json:"bookId" db:"book_id"`
	Title         string     `json:"title" db:"title"`
	Content       *string    `json:"content,omitempty" db:"content"`
	ChapterOrder  int        `json:"chapterOrder" db:"chapter_order"`
	Status        string     `json:"status" db:"status"`
	CoinCost      int        `json:"coinCost" db:"coin_cost"`
	TotalViews    int64      `json:"totalViews" db:"total_views"`
	CreateDatetime time.Time `json:"createDatetime" db:"create_datetime"`
	UpdateDatetime *time.Time `json:"updateDatetime,omitempty" db:"update_datetime"`
}

// --- STRUCT BARU UNTUK RESPONSE DETAIL ---
// BookDetailResponse adalah struktur data yang menyeluruh untuk halaman detail buku.
type BookDetailResponse struct {
    BookInfo    *Book      `json:"bookInfo"`
    Chapters    []Chapter  `json:"chapters"`
    Author      *User      `json:"author"`
    // Anda bisa menambahkan data lain di sini nanti, seperti total ulasan, dll.
}
