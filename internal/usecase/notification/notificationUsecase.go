package notification

import (
	"fmt"
	"notification/internal/entity"
	log "notification/internal/platform/repositories"
	"notification/internal/usecase/notifiers"
	"time"
)

type NotificationUseCase struct {
	LogRepository log.Log
	SMSUsecase    *notifiers.SMSUsecase
	EmailUsecase  *notifiers.EmailUsecase
	PushUsecase   *notifiers.PushUsecase
}

type Notification interface {
	SendNotification(user entity.User, message string) error
}

func NewNotificationUseCase(log log.Log) *NotificationUseCase {
	smsUsecase := &notifiers.SMSUsecase{}
	emailUsecase := &notifiers.EmailUsecase{}
	pushUsecase := &notifiers.PushUsecase{}

	return &NotificationUseCase{
		LogRepository: log,
		SMSUsecase:    smsUsecase,
		EmailUsecase:  emailUsecase,
		PushUsecase:   pushUsecase,
	}
}

func (n NotificationUseCase) SendNotification(notification entity.Notification) ([]entity.Log, error) {
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

func (n NotificationUseCase) GetUsersByCategory(category entity.Category) []entity.User {
	var users []entity.User

	switch category {
	case entity.SportsCategory:
		users = []entity.User{
			{
				ID:          1,
				Name:        "Mary Alexander",
				Email:       "mary.alexander@outlook.com",
				PhoneNumber: "78958745",
				Subscribed:  []entity.Category{entity.SportsCategory},
				Channels:    []entity.Channel{"SMS"},
			},
		}
	case entity.FinanceCategory:
		users = []entity.User{
			{
				ID:          2,
				Name:        "Antony Smith",
				Email:       "antony.smith@gmail.com",
				PhoneNumber: "4134132441",
				Subscribed:  []entity.Category{entity.FinanceCategory},
				Channels:    []entity.Channel{"Email", "Push"},
			},
		}
	case entity.MoviesCategory:
		users = []entity.User{
			{
				ID:          3,
				Name:        "Any Johnson",
				Email:       "any.johnson@gmail.com",
				PhoneNumber: "+123456789",
				Subscribed:  []entity.Category{entity.MoviesCategory},
				Channels:    []entity.Channel{"SMS", "Email"},
			},
			{
				ID:          4,
				Name:        "Fred Williams",
				Email:       "fred.williams@hotmail.com",
				PhoneNumber: "78459214465",
				Subscribed:  []entity.Category{entity.MoviesCategory},
				Channels:    []entity.Channel{"Email"},
			},
		}
	}

	return users
}

func (n NotificationUseCase) GetLogs() ([]entity.Log, error) {
	return n.LogRepository.GetLogs()
}

func (n NotificationUseCase) DeleteLogs() error {
	return n.LogRepository.DeleteLogs()
}

func (n NotificationUseCase) send(notification entity.Notification, user entity.User) ([]entity.Log, error) {
	var logs []entity.Log
	notifiers := n.getNotifiers(user.Channels)
	for _, notifier := range notifiers {
		err := notifier.SendNotification(user, notification.Message)
		if err != nil {
			return nil, err
		}

		notificationType := n.getNotificationType(notifier)

		log := entity.Log{
			ID:               fmt.Sprintf("%v-%s-%s", user.ID, notification.Category, notificationType),
			UserID:           user.ID,
			Message:          notification.Message,
			Category:         notification.Category,
			NotificationType: notificationType,
			Timestamp:        time.Now(),
		}

		n.LogRepository.SaveLog(log)
		logs = append(logs, log)
	}

	return logs, nil
}

func (n NotificationUseCase) getNotifiers(channels []entity.Channel) []Notification {
	notifiers := make([]Notification, 0)
	fmt.Println(len(notifiers))

	for _, channel := range channels {
		switch channel {
		case entity.EmailChannel:
			notifiers = append(notifiers, n.EmailUsecase)
		case entity.SMSChannel:
			notifiers = append(notifiers, n.SMSUsecase)
		case entity.PushChannel:
			notifiers = append(notifiers, n.PushUsecase)
		}

	}

	return notifiers
}

func (n NotificationUseCase) getNotificationType(notifier Notification) string {
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
