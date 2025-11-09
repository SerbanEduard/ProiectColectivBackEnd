package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFriendRequestController_SendFriendRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockFriendRequestService)
	controller := controller.NewFriendRequestController()
	controller.SetFriendRequestService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "fromUserId", Value: "user1"},
		{Key: "toUserId", Value: "user2"},
	}

	mockService.On("SendFriendRequest", "user1", "user2").Return(nil)

	controller.SendFriendRequest(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestFriendRequestController_RespondToRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockFriendRequestService)
	controller := controller.NewFriendRequestController()
	controller.SetFriendRequestService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "fromUserId", Value: "user1"},
		{Key: "toUserId", Value: "user2"},
	}

	requestBody := map[string]bool{"accept": true}
	jsonData, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest("PUT", "/friend-requests/user1/user2", bytes.NewBuffer(jsonData))

	mockService.On("RespondToFriendRequest", "user1", "user2", true).Return(nil)

	controller.RespondToFriendRequest(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestFriendRequestController_SendFriendRequest_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockFriendRequestService)
	controller := controller.NewFriendRequestController()
	controller.SetFriendRequestService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "fromUserId", Value: "invalidUser"},
		{Key: "toUserId", Value: "user2"},
	}

	mockService.On("SendFriendRequest", "invalidUser", "user2").Return(fmt.Errorf("sender user not found"))

	controller.SendFriendRequest(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "sender user not found")
}

func TestFriendRequestController_RespondToRequest_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockFriendRequestService)
	controller := controller.NewFriendRequestController()
	controller.SetFriendRequestService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "fromUserId", Value: "user1"},
		{Key: "toUserId", Value: "user2"},
	}

	c.Request, _ = http.NewRequest("PUT", "/friend-requests/user1/user2", bytes.NewBuffer([]byte(`{invalid json}`)))

	controller.RespondToFriendRequest(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFriendRequestController_SendFriendRequest_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockFriendRequestService)
	controller := controller.NewFriendRequestController()
	controller.SetFriendRequestService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "fromUserId", Value: ""},
		{Key: "toUserId", Value: "user2"},
	}

	controller.SendFriendRequest(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid user IDs")
}

func TestFriendRequestController_RespondToRequest_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockFriendRequestService)
	controller := controller.NewFriendRequestController()
	controller.SetFriendRequestService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "fromUserId", Value: "user1"},
		{Key: "toUserId", Value: ""},
	}

	controller.RespondToFriendRequest(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid user IDs")
}

func TestFriendRequestController_GetPendingRequests_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockFriendRequestService)
	controller := controller.NewFriendRequestController()
	controller.SetFriendRequestService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "userId", Value: "invalidUser"},
	}

	mockService.On("GetPendingRequests", "invalidUser").Return(nil, fmt.Errorf("user not found"))

	controller.GetPendingRequests(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
