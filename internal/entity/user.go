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
	SportsCategory  Category = "Sports"
	FinanceCategory Category = "Finance"
	MoviesCategory  Category = "Movies"
)

type Channel string

const (
	EmailChannel Channel = "Email"
	PushChannel  Channel = "Push"
	SMSChannel   Channel = "SMS"
)
