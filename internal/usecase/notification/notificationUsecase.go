package notification

import (
	"notification/internal/entity"
	log "notification/internal/platform/repositories"
	email "notification/internal/usecase/email"
	push "notification/internal/usecase/push"
	sms "notification/internal/usecase/sms"
	"time"
)

type NotificationUseCase struct {
	LogRepository           *log.LogRepository
	SMSService              *sms.SMSUsecase
	EmailService            *email.EmailUsecase
	PushNotificationService *push.PushUsecase
	Observers               []entity.Observer
}

func NewNotificationUseCase() *NotificationUseCase {
	logRepository := &log.LogRepository{}
	smsService := &sms.SMSUsecase{}
	emailService := &email.EmailUsecase{}
	pushNotificationService := &push.PushUsecase{}
	observers := make([]entity.Observer, 0)

	return &NotificationUseCase{
		LogRepository:           logRepository,
		SMSService:              smsService,
		EmailService:            emailService,
		PushNotificationService: pushNotificationService,
		Observers:               observers,
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

func (app *NotificationUseCase) GetNotifiers(category entity.Category) []entity.Notifier {
	switch category {
	case entity.Sports:
		return []entity.Notifier{app.SMSService}
	case entity.Finance:
		return []entity.Notifier{app.EmailService, app.PushNotificationService}
	case entity.Movies:
		return []entity.Notifier{app.EmailService}
	default:
		return []entity.Notifier{}
	}
}

func (app *NotificationUseCase) SendNotification(notification entity.Notification) ([]entity.Log, error) {
	var logs []entity.Log

	for _, notifier := range notification.Notifiers {
		users := app.GetUsersByCategory(notification.Category)
		for _, user := range users {
			err := notifier.SendNotification(user, notification.Message)
			if err != nil {
				// TODO: Handle error
			}

			log := entity.Log{
				Message:          notification.Message,
				Category:         notification.Category,
				NotificationType: GetNotifierType(notifier),
				Timestamp:        time.Now(),
			}

			app.LogRepository.SaveLog(log)
			logs = append(logs, log)
		}
	}

	return logs, nil
}

// func (app *NotificationUseCase) Send(notifier entity.Notifier, user entity.User) ([]entity.Log, error) {
// 	for _, channel := range user.Channels {

// 	}

// }

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
				Channels:    []string{"SMS", "E-Mail", "Push Notification"},
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
				Channels:    []string{"E-Mail", "Push Notification"},
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
				Channels:    []string{"SMS", "E-Mail"},
			},
			{
				ID:          4,
				Name:        "Fred Williams",
				Email:       "fred.williams@hotmail.com",
				PhoneNumber: "78459214465",
				Subscribed:  []entity.Category{entity.Movies},
				Channels:    []string{"E-Mail"},
			},
		}
	}

	return users
}

func (app *NotificationUseCase) GetLogs() []entity.Log {
	return app.LogRepository.GetLogs()
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
