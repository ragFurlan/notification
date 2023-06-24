package entity

import "time"

type Log struct {
	ID               int
	UserID           int
	Message          string
	Category         Category
	NotificationType string
	Timestamp        time.Time
}
