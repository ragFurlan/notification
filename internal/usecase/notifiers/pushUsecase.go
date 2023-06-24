package notifiers

import (
	"fmt"
	"notification/internal/entity"
)

type PushUsecase struct{}

func (s *PushUsecase) SendNotification(user entity.User, message string) error {
	fmt.Printf("Sending push notification to %s: %s\n", user.Name, message)
	// TODO:Logic to send push notification
	return nil
}
