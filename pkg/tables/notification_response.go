package tables

import "time"

// NotificationResponse adalah struct gabungan untuk response API notifikasi.
type NotificationResponse struct {
	NotificationID   int64     `json:"notificationId" db:"notification_id"`
	NotificationType string    `json:"notificationType" db:"notification_type"`
	IsRead           bool      `json:"isRead" db:"is_read"`
	CreateDatetime   time.Time `json:"createDatetime" db:"create_datetime"`

	ActorName      *string `json:"userName,omitempty" db:"actor_name"`
	ActorAvatarURL *string `json:"userAvatarUrl,omitempty" db:"actor_avatar_url"`
	
	BookName       *string `json:"bookName,omitempty" db:"book_name"`
	ChapterName    *string `json:"chapterName,omitempty" db:"chapter_name"`
	CommentContent *string `json:"commentContent,omitempty" db:"comment_content"`

	// ✨ ID BARU UNTUK NAVIGASI
	BookID           *int64 `json:"bookId,omitempty" db:"book_id"`
	RelatedEntityID  *int64 `json:"relatedEntityId,omitempty" db:"related_entity_id"`
}

// PaginatedNotificationResponse adalah struktur untuk response daftar notifikasi yang dibungkus.
type PaginatedNotificationResponse struct {
    Pagination PaginationInfo         `json:"pagination"`
    Items      []NotificationResponse `json:"notifications"`
}