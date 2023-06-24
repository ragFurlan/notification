package entity

type User struct {
	ID          int
	Name        string
	Email       string
	PhoneNumber string
	Subscribed  []Category
	Channels    []Channel
}

type Category string

const (
	Sports  Category = "Sports"
	Finance Category = "Finance"
	Movies  Category = "Movies"
)

type Channel string

const (
	Email Channel = "Email"
	Push  Channel = "Push"
	SMS   Channel = "SMS"
)

func (u *User) Update(notification Notification) {
	//  TODO: Handle the received notification for the user
	// You can implement the specific logic here, such as sending the notification through the user's preferred channels (SMS, email, etc.)
}
