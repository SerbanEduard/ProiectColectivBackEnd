package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserController_GetFriends_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(tests.MockFriendRequestService)
	uc := controller.NewUserController()
	uc.SetFriendRequestService(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "user1"}}

	mockSvc.On("GetFriends", "user1").Return([]*entity.User{{ID: "f1", Username: "u1"}}, nil)

	uc.GetFriends(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []*entity.User
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	mockSvc.AssertExpectations(t)
}

func TestUserController_GetMutualFriends_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(tests.MockFriendRequestService)
	uc := controller.NewUserController()
	uc.SetFriendRequestService(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: "userA"},
		{Key: "otherId", Value: "userB"},
	}

	mockSvc.On("GetMutualFriends", "userA", "userB").Return([]*entity.User{{ID: "m1", Username: "mu"}}, nil)

	uc.GetMutualFriends(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []*entity.User
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	mockSvc.AssertExpectations(t)
}
