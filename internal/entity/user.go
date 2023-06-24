package entity

type User struct {
	ID          int
	Name        string
	Email       string
	PhoneNumber string
	Subscribed  []Category
	Channels    []string
}

type Category string

const (
	Sports  Category = "Sports"
	Finance Category = "Finance"
	Movies  Category = "Movies"
)

func (u *User) Update(notification Notification) {
	//  TODO: Handle the received notification for the user
	// You can implement the specific logic here, such as sending the notification through the user's preferred channels (SMS, email, etc.)
}
