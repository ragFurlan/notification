package entity

type Observer interface {
	Update(notification Notification)
}
