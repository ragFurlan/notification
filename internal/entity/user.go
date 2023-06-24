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
