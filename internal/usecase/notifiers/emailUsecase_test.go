package notifiers

import (
	"notification/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail_Success(t *testing.T) {
	service := EmailUsecase{}
	err := service.SendNotification(entity.User{}, "test function")
	assert.NoError(t, err)

}
