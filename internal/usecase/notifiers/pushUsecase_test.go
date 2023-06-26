package notifiers

import (
	"notification/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush_Success(t *testing.T) {
	service := PushUsecase{}
	err := service.SendNotification(entity.User{}, "test function")
	assert.NoError(t, err)

}
