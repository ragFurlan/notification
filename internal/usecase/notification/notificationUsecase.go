package notification

import (
	"fmt"
	"notification/internal/entity"
	log "notification/internal/platform/repositories"
	email "notification/internal/usecase/email"
	push "notification/internal/usecase/push"
	sms "notification/internal/usecase/sms"
	"time"

	"github.com/google/uuid"
)

type NotificationUseCase struct {
	LogRepository *log.LogRepository
	SMSUsecase    *sms.SMSUsecase
	EmailUsecase  *email.EmailUsecase
	PushUsecase   *push.PushUsecase
	Observers     []entity.Observer
}

func NewNotificationUseCase(logRepository log.LogRepository) *NotificationUseCase {

	smsUsecase := &sms.SMSUsecase{}
	emailUsecase := &email.EmailUsecase{}
	pushUsecase := &push.PushUsecase{}
	observers := make([]entity.Observer, 0)

	return &NotificationUseCase{
		LogRepository: &logRepository,
		SMSUsecase:    smsUsecase,
		EmailUsecase:  emailUsecase,
		PushUsecase:   pushUsecase,
		Observers:     observers,
	}
}

func (app *NotificationUseCase) RegisterObserver(observer entity.Observer) {
	app.Observers = append(app.Observers, observer)
}

func (app *NotificationUseCase) NotifyObservers(notification entity.Notification) {
	for _, observer := range app.Observers {
		observer.Update(notification)
	}
}

func (app *NotificationUseCase) GetNotifiers(channels []entity.Channel) []entity.Notifier {
	notifiers := make([]entity.Notifier, 0)
	fmt.Println(len(notifiers))

	for _, channel := range channels {
		switch channel {
		case entity.Email:
			notifiers = append(notifiers, app.EmailUsecase)
		case entity.SMS:
			notifiers = append(notifiers, app.SMSUsecase)
		case entity.Push:
			notifiers = append(notifiers, app.PushUsecase)
		}

	}

	return notifiers
}

func (app *NotificationUseCase) SendNotification(notification entity.Notification) ([]entity.Log, error) {
	var logs []entity.Log
	users := app.GetUsersByCategory(notification.Category)
	for _, user := range users {
		logsOfUsers, err := app.Send(notification, user)
		if err != nil {
			// TODO: Handle error
		}

		for _, log := range logsOfUsers {
			logs = append(logs, log)
		}

	}

	return logs, nil

}

func (app *NotificationUseCase) Send(notification entity.Notification, user entity.User) ([]entity.Log, error) {
	var logs []entity.Log
	notifiers := app.GetNotifiers(user.Channels)
	for _, notifier := range notifiers {
		err := notifier.SendNotification(user, notification.Message)
		if err != nil {
			// TODO: Handle error
		}

		log := entity.Log{
			ID:               uuid.New().String(),
			UserID:           user.ID,
			Message:          notification.Message,
			Category:         notification.Category,
			NotificationType: GetNotifierType(notifier),
			Timestamp:        time.Now(),
		}

		app.LogRepository.SaveLog(log)
		logs = append(logs, log)
	}

	return logs, nil
}

func (app *NotificationUseCase) GetUsersByCategory(category entity.Category) []entity.User {
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

func (app *NotificationUseCase) GetLogs() ([]entity.Log, error) {
	return app.LogRepository.GetLogs()
}

func (app *NotificationUseCase) DeleteLogs() error {
	return app.LogRepository.DeleteLogs()
}

// Helper function to get the notifier type as a string
func GetNotifierType(notifier entity.Notifier) string {
	switch notifier.(type) {
	case *sms.SMSUsecase:
		return "SMS"
	case *email.EmailUsecase:
		return "E-Mail"
	case *push.PushUsecase:
		return "Push Notification"
	default:
		return "Unknown"
	}
}
