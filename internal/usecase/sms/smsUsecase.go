package sms

import (
	"fmt"
	"notification/internal/entity"
)

type SMSUsecase struct{}

func (s *SMSUsecase) SendNotification(user entity.User, message string) error {
	fmt.Printf("Sending SMS notification to %s (%s): %s\n", user.Name, user.PhoneNumber, message)
	// TODO: Logic to send SMS notification
	return nil
}
