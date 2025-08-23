package dao

import (
	"context"
	"fmt"
	"noversystem/pkg/tables"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationDao struct {
	DB *pgxpool.Pool
}

func NewNotificationDao(db *pgxpool.Pool) *NotificationDao {
	return &NotificationDao{DB: db}
}

// GetNotificationsByUserID mengambil daftar notifikasi untuk pengguna dengan pagination.
func (d *NotificationDao) GetNotificationsByUserID(ctx context.Context, userID int64, limit, offset int) ([]tables.NotificationResponse, error) {
	var notifications []tables.NotificationResponse
	
	// ✨ QUERY DEFINITIF YANG SUDAH DIPERBAIKI TOTAL
	query := `
		SELECT
			sn.notification_id,
			sn.is_read,
			sn.create_datetime,
			sn.notification_type,
			sn.related_entity_id,
			
			COALESCE(actor.pen_name, actor.full_name, 'Anonymous') AS actor_name,
			actor.avatar_url AS actor_avatar_url,
			
			COALESCE(b1.title, b2.title, b3.title, b4.title) AS book_name,
			COALESCE(b1.book_id, b2.book_id, b3.book_id, b4.book_id) AS book_id,
			COALESCE(c1.title, c2.title) AS chapter_name,
			COALESCE(bc.comment_text, cc.comment_text) AS comment_content
		FROM 
			system_notifications sn
		LEFT JOIN 
			users actor ON sn.actor_id = actor.user_id
		LEFT JOIN 
			book_comments bc ON sn.related_entity_id = bc.comment_id AND sn.related_entity_type = 'BOOK_COMMENT' -- Lebih spesifik
		LEFT JOIN 
			books b1 ON bc.book_id = b1.book_id
		LEFT JOIN 
			chapter_comments cc ON sn.related_entity_id = cc.comment_id AND sn.related_entity_type = 'CHAPTER_COMMENT' -- Lebih spesifik
		LEFT JOIN 
			chapters c1 ON cc.chapter_id = c1.chapter_id
		LEFT JOIN 
			books b2 ON c1.book_id = b2.book_id
		LEFT JOIN 
			chapters c2 ON sn.related_entity_id = c2.chapter_id AND sn.related_entity_type = 'CHAPTER'
		LEFT JOIN 
			books b3 ON c2.book_id = b3.book_id
		LEFT JOIN
			books b4 ON sn.related_entity_id = b4.book_id AND sn.related_entity_type = 'BOOK'
		WHERE 
			sn.user_id = $1
		ORDER BY 
			sn.create_datetime DESC
		LIMIT $2 OFFSET $3`

	err := pgxscan.Select(ctx, d.DB, &notifications, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil notifikasi: %w", err)
	}
	return notifications, nil
}

// CountNotificationsByUserID menghitung total notifikasi untuk seorang pengguna.
func (d *NotificationDao) CountNotificationsByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM system_notifications WHERE user_id = $1`
	err := d.DB.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}