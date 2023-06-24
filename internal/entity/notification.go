package entity

type Notification struct {
	Message   string
	Category  Category
	Notifiers []Notifier
}

type Notifier interface {
	SendNotification(user User, message string) error
}
