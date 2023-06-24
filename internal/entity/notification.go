package entity

type Notification struct {
	Message  string
	Category Category
}

type Notifier interface {
	SendNotification(user User, message string) error
}
