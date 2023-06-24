package email

import (
	"fmt"
	"notification/internal/entity"
)

type EmailUsecase struct{}

func (s *EmailUsecase) SendNotification(user entity.User, message string) error {
	fmt.Printf("Sending email notification to %s (%s): %s\n", user.Name, user.Email, message)
	// Logic to send email notification
	return nil
}
