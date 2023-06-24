package notification

import (
	"fmt"
	"notification/internal/entity"
	log "notification/internal/platform/repositories"
	"notification/internal/usecase/notifiers"
	"time"

	"github.com/google/uuid"
)

type NotificationUseCase struct {
	LogRepository log.Log
	SMSUsecase    *notifiers.SMSUsecase
	EmailUsecase  *notifiers.EmailUsecase
	PushUsecase   *notifiers.PushUsecase
	Observers     []entity.Observer
}

type Notifier interface {
	SendNotification(user entity.User, message string) error
}

func NewNotificationUseCase(log log.Log) *NotificationUseCase {
	smsUsecase := &notifiers.SMSUsecase{}
	emailUsecase := &notifiers.EmailUsecase{}
	pushUsecase := &notifiers.PushUsecase{}
	observers := make([]entity.Observer, 0)

	return &NotificationUseCase{
		LogRepository: log,
		SMSUsecase:    smsUsecase,
		EmailUsecase:  emailUsecase,
		PushUsecase:   pushUsecase,
		Observers:     observers,
	}
}

func (n *NotificationUseCase) SendNotification(notification entity.Notification) ([]entity.Log, error) {
	var logs []entity.Log
	users := n.GetUsersByCategory(notification.Category)
	for _, user := range users {
		logsOfUsers, err := n.send(notification, user)
		if err != nil {
			return nil, err
		}

		for _, log := range logsOfUsers {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

func (n *NotificationUseCase) GetUsersByCategory(category entity.Category) []entity.User {
	var users []entity.User

	switch category {
	case entity.Sports:
		users = []entity.User{
			{
				ID:          1,
				Name:        "Mary Alexander",
				Email:       "mary.alexander@outlook.com",
				PhoneNumber: "78958745",
				Subscribed:  []entity.Category{entity.Sports},
				Channels:    []entity.Channel{"SMS", "Email", "Push"},
			},
		}
	case entity.Finance:
		users = []entity.User{
			{
				ID:          2,
				Name:        "Antony Smith",
				Email:       "antony.smith@gmail.com",
				PhoneNumber: "4134132441",
				Subscribed:  []entity.Category{entity.Finance},
				Channels:    []entity.Channel{"Email", "Push"},
			},
		}
	case entity.Movies:
		users = []entity.User{
			{
				ID:          3,
				Name:        "Any Johnson",
				Email:       "any.johnson@gmail.com",
				PhoneNumber: "+123456789",
				Subscribed:  []entity.Category{entity.Movies},
				Channels:    []entity.Channel{"SMS", "Email"},
			},
			{
				ID:          4,
				Name:        "Fred Williams",
				Email:       "fred.williams@hotmail.com",
				PhoneNumber: "78459214465",
				Subscribed:  []entity.Category{entity.Movies},
				Channels:    []entity.Channel{"Email"},
			},
		}
	}

	return users
}

func (n *NotificationUseCase) GetLogs() ([]entity.Log, error) {
	return n.LogRepository.GetLogs()
}

func (n *NotificationUseCase) DeleteLogs() error {
	return n.LogRepository.DeleteLogs()
}

func (n *NotificationUseCase) send(notification entity.Notification, user entity.User) ([]entity.Log, error) {
	var logs []entity.Log
	notifiers := n.getNotifiers(user.Channels)
	for _, notifier := range notifiers {
		err := notifier.SendNotification(user, notification.Message)
		if err != nil {
			return nil, err
		}

		log := entity.Log{
			ID:               uuid.New().String(),
			UserID:           user.ID,
			Message:          notification.Message,
			Category:         notification.Category,
			NotificationType: n.getNotifierType(notifier),
			Timestamp:        time.Now(),
		}

		n.LogRepository.SaveLog(log)
		logs = append(logs, log)
	}

	return logs, nil
}

func (n *NotificationUseCase) getNotifiers(channels []entity.Channel) []Notifier {
	notifiers := make([]Notifier, 0)
	fmt.Println(len(notifiers))

	for _, channel := range channels {
		switch channel {
		case entity.Email:
			notifiers = append(notifiers, n.EmailUsecase)
		case entity.SMS:
			notifiers = append(notifiers, n.SMSUsecase)
		case entity.Push:
			notifiers = append(notifiers, n.PushUsecase)
		}

	}

	return notifiers
}

func (n *NotificationUseCase) getNotifierType(notifier Notifier) string {
	switch notifier.(type) {
	case *notifiers.SMSUsecase:
		return "SMS"
	case *notifiers.EmailUsecase:
		return "E-Mail"
	case *notifiers.PushUsecase:
		return "Push Notification"
	default:
		return "Unknown"
	}
}
