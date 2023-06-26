package notifiers

import (
	"notification/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMS_Success(t *testing.T) {
	service := SMSUsecase{}
	err := service.SendNotification(entity.User{}, "test function")
	assert.NoError(t, err)

}
