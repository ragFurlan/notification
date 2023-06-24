package notification

import (
	"errors"
	"fmt"
	"notification/internal/entity"
	log "notification/test/platform"
	//notification "notification/test/usecase"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var anyError = errors.New("Error")

func TestNotification_SendNotification_Success(t *testing.T) {
	controller := gomock.NewController(t)

	defer controller.Finish()
	logEntity := log.NewMockLog(controller)
	service := NewNotificationUseCase(logEntity)

	message := getMessage(1, "SMS")
	logEntity.EXPECT().SaveLog(message).Return(nil)

	_, err := service.SendNotification(getNotification())
	assert.NoError(t, err)

}

func TestSendNotification_GetLogs_Success(t *testing.T) {
	controller := gomock.NewController(t)

	defer controller.Finish()
	logEntity := log.NewMockLog(controller)
	service := NewNotificationUseCase(logEntity)

	logEntity.EXPECT().GetLogs().Return([]entity.Log{}, nil)

	_, err := service.GetLogs()
	assert.NoError(t, err)

}

func TestSendNotification_DeleteLogs_Success(t *testing.T) {
	controller := gomock.NewController(t)

	defer controller.Finish()
	logEntity := log.NewMockLog(controller)
	service := NewNotificationUseCase(logEntity)

	logEntity.EXPECT().DeleteLogs().Return(nil)

	err := service.DeleteLogs()
	assert.NoError(t, err)

}

func getUser(id int) entity.User {
	return entity.User{
		ID:          id,
		Name:        "Mary Alexander",
		Email:       "mary.alexander@outlook.com",
		PhoneNumber: "78958745",
		Subscribed:  []entity.Category{entity.SportsCategory},
		Channels:    []entity.Channel{"SMS"},
	}
}

func getMessage(id int, NotificationType string) entity.Log {
	return entity.Log{
		ID:               fmt.Sprintf("%v-%s-%s", id, string(entity.SportsCategory), string(NotificationType)),
		UserID:           id,
		Message:          "test test",
		Category:         entity.SportsCategory,
		NotificationType: NotificationType,
		Timestamp:        time.Now(),
	}
}

func getNotification() entity.Notification {
	return entity.Notification{
		Message:  "test test",
		Category: entity.SportsCategory,
	}
}
