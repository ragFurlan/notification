package log

import (
	"fmt"
	"notification/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

)

func TestLog_Success(t *testing.T) {
	urlLog := "./logs.txt"

	logRepository := NewLogRepository(urlLog)

	log := getMessage(1, "SMS")
	err := logRepository.SaveLog(log)
	assert.NoError(t, err)

	logs, err := logRepository.GetLogs()
	assert.Equal(t, len(logs), 1)

	err = logRepository.DeleteLogs()
	assert.NoError(t, err)

	_, err = logRepository.GetLogs()
	assert.Error(t, err)

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
