package tables

import "time"

// SystemNotification merepresentasikan record dalam tabel system_notifications.
type SystemNotification struct {
	NotificationID    int64      `json:"notificationId"`
	UserID            int64      `json:"userId"`
	ActorID           *int64     `json:"actorId,omitempty"` // ✨ DITAMBAHKAN
	NotificationType  string     `json:"notificationType"`
	Content           string     `json:"content"`
	IsRead            bool       `json:"isRead"`
	RelatedEntityType *string    `json:"relatedEntityType,omitempty"`
	RelatedEntityID   *int64     `json:"relatedEntityId,omitempty"`
	CreateDatetime    time.Time  `json:"createDatetime"`
	UpdateDatetime    *time.Time `json:"updateDatetime,omitempty"`
}