package notification_handler

import (
	//"fmt"
	//"bytes"
	//"encoding/json"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	//"encoding/json"
	//"io"
	"net/http"
	"net/http/httptest"
	"notification/internal/entity"
	usecase "notification/internal/usecase/notification"
	log "notification/test/platform"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	handler     *NotificationHandler
	logMock     *log.MockLog
	controller  *gomock.Controller
	usecaseMock *usecase.NotificationUseCase
	anyError    = errors.New("Error")
)

func (mock *MockHTTP) Do(_ *http.Request) (*http.Response, error) {
	return mock.response, mock.err
}

type MockHTTP struct {
	response *http.Response
	err      error
}

func TestSubmitNotification_Success(t *testing.T) {
	bodyReader := strings.NewReader(`{"category": "Sports", "message": "Test Submit Notification"}`)
	r := httptest.NewRequest(http.MethodPost, "/add", bodyReader)
	w := httptest.NewRecorder()
	setHandlerAndLogMock(t)

	message := getMessage(1, "SMS")
	logMock.EXPECT().SaveLog(message).Return(nil)

	handler.SubmitNotification(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusOK, got.StatusCode)
	controller.Finish()
}

func TestSubmitNotification_Method_Error(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/add", nil)
	w := httptest.NewRecorder()

	setHandlerAndLogMock(t)
	handler.SubmitNotification(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, got.StatusCode)
	controller.Finish()
}

func TestSubmitNotification_Error(t *testing.T) {
	bodyReader := strings.NewReader(`{"category": "Sports", "message": "Test Submit Notification"}`)
	r := httptest.NewRequest(http.MethodPost, "/add", bodyReader)
	w := httptest.NewRecorder()
	setHandlerAndLogMock(t)

	message := getMessage(1, "SMS")
	logMock.EXPECT().SaveLog(message).Return(anyError)

	handler.SubmitNotification(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusInternalServerError, got.StatusCode)
	controller.Finish()
}

func TestSubmitNotification_Body_Success(t *testing.T) {
	bodyReader := strings.NewReader(`[{}]`)
	r := httptest.NewRequest(http.MethodPost, "/add", bodyReader)
	w := httptest.NewRecorder()
	setHandlerAndLogMock(t)
	handler.SubmitNotification(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusBadRequest, got.StatusCode)
	controller.Finish()
}

func TestGetLogsAndDeleteLogs_Success(t *testing.T) {
	bodyReader := strings.NewReader(`{"category": "Sports", "message": "Test Submit Notification"}`)
	r := httptest.NewRequest(http.MethodGet, "/get", bodyReader)
	w := httptest.NewRecorder()
	setHandlerAndLogMock(t)

	// GET
	message := getMessage(1, "SMS")
	logMock.EXPECT().GetLogs().Return([]entity.Log{message}, nil)
	handler.GetLogs(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusOK, got.StatusCode)

	var requestBody []entity.Log
	err := json.NewDecoder(got.Body).Decode(&requestBody)
	assert.NoError(t, err)

	assert.Equal(t, requestBody[0].Category, entity.SportsCategory)
	assert.Equal(t, requestBody[0].Message, "Test Submit Notification")

	// DELETE
	r = httptest.NewRequest(http.MethodDelete, "/delete", nil)
	w = httptest.NewRecorder()

	logMock.EXPECT().DeleteLogs().Return(nil)
	handler.DeleteLogs(w, r)

	got = w.Result()
	assert.Equal(t, http.StatusOK, got.StatusCode)
	controller.Finish()
}

func TestGetLogs_Method_Error(t *testing.T) {
	r := httptest.NewRequest(http.MethodDelete, "/get", nil)
	w := httptest.NewRecorder()

	setHandlerAndLogMock(t)
	handler.GetLogs(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, got.StatusCode)
	controller.Finish()
}

func TestGetLogs_Error(t *testing.T) {
	bodyReader := strings.NewReader(`{"category": "Sports", "message": "Test Submit Notification"}`)
	r := httptest.NewRequest(http.MethodGet, "/get", bodyReader)
	w := httptest.NewRecorder()
	setHandlerAndLogMock(t)

	// GET
	message := getMessage(1, "SMS")
	logMock.EXPECT().GetLogs().Return([]entity.Log{message}, anyError)
	handler.GetLogs(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusInternalServerError, got.StatusCode)

}

func TestDeleteLogs_Method_Error(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/delete", nil)
	w := httptest.NewRecorder()

	setHandlerAndLogMock(t)
	handler.DeleteLogs(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, got.StatusCode)
	controller.Finish()
}

func TestDeleteLogs_Error(t *testing.T) {
	bodyReader := strings.NewReader(`{"category": "Sports", "message": "Test Submit Notification"}`)
	r := httptest.NewRequest(http.MethodGet, "/get", bodyReader)
	w := httptest.NewRecorder()
	setHandlerAndLogMock(t)

	// GET
	message := getMessage(1, "SMS")
	logMock.EXPECT().GetLogs().Return([]entity.Log{message}, nil)
	handler.GetLogs(w, r)

	got := w.Result()
	assert.Equal(t, http.StatusOK, got.StatusCode)

	var requestBody []entity.Log
	err := json.NewDecoder(got.Body).Decode(&requestBody)
	assert.NoError(t, err)

	assert.Equal(t, requestBody[0].Category, entity.SportsCategory)
	assert.Equal(t, requestBody[0].Message, "Test Submit Notification")

	// DELETE
	r = httptest.NewRequest(http.MethodDelete, "/delete", nil)
	w = httptest.NewRecorder()

	logMock.EXPECT().DeleteLogs().Return(anyError)
	handler.DeleteLogs(w, r)

	got = w.Result()
	assert.Equal(t, http.StatusInternalServerError, got.StatusCode)
	controller.Finish()
}

func setHandlerAndLogMock(t *testing.T) {
	controller = gomock.NewController(t)
	logMock = log.NewMockLog(controller)
	usecaseMock = usecase.NewNotificationUseCase(logMock)
	handler = NewNotificationHandler(usecaseMock)
}

func getMessage(id int, NotificationType string) entity.Log {
	return entity.Log{
		ID:               fmt.Sprintf("%v-%s-%s", id, string(entity.SportsCategory), string(NotificationType)),
		UserID:           id,
		Message:          "Test Submit Notification",
		Category:         entity.SportsCategory,
		NotificationType: NotificationType,
		Timestamp:        time.Now(),
	}
}
