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
    Genres        *string    `json:"genres,omitempty" db:"genres"` 
    AuthorID      int64      `json:"-" db:"author_id"`
    AuthorPenName *string    `json:"authorPenName,omitempty" db:"pen_name"`
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

// --- STRUCT BARU UNTUK Ulasan ---
// Review merepresentasikan satu ulasan dari pengguna untuk sebuah buku.
type Review struct {
	ReviewID       int64      `json:"reviewId" db:"review_id"`
	UserID         int64      `json:"userId" db:"user_id"`
	AuthorPenName  *string    `json:"authorPenName,omitempty" db:"pen_name"`
	AuthorAvatar   *string    `json:"authorAvatar,omitempty" db:"avatar_url"`
	Rating         int        `json:"rating" db:"rating"`
	ReviewText     *string    `json:"reviewText,omitempty" db:"review_text"`
	CreateDatetime time.Time  `json:"createDatetime" db:"create_datetime"`
}


// BookDetailResponse adalah struktur data yang menyeluruh untuk halaman detail buku.
type BookDetailResponse struct {
    BookInfo    *Book      `json:"bookInfo"`
    Chapters    []Chapter  `json:"chapters"`
    Author      *User      `json:"author"`
    Reviews     []Review   `json:"reviews"` // Ditambahkan untuk menampung ulasan
}

// PaginatedBookResponse adalah struktur untuk response daftar buku yang disertai info pagination.
type PaginatedBookResponse struct {
    Pagination PaginationInfo `json:"pagination"`
    Books      []Book         `json:"books"`
}

// PaginationInfo berisi detail tentang pagination.
type PaginationInfo struct {
    CurrentPage int   `json:"currentPage"`
    PageSize    int   `json:"pageSize"`
    TotalItems  int64 `json:"totalItems"`
    TotalPages  int   `json:"totalPages"`
}
