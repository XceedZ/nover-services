package tables

import "time"

// ChapterComment merepresentasikan record dalam tabel chapter_comments.
type ChapterComment struct {
	CommentID       int64      `json:"commentId"`
	ChapterID       int64      `json:"chapterId"`
	UserID          int64      `json:"userId"`
	CommentText     string     `json:"commentText"`
	ParentCommentID *int64     `json:"parentCommentId,omitempty"` // Pointer untuk nilai NULL
	CreateDatetime  time.Time  `json:"createDatetime"`
	UpdateDatetime  *time.Time `json:"updateDatetime,omitempty"`   // Pointer untuk nilai NULL
}