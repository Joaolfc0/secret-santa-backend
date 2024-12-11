package group

import (
	"net/http"
	"testing"

	"service-secret-santa/customError"
	"service-secret-santa/functions"
	"service-secret-santa/models"
	mocks "service-secret-santa/services/group/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func internalErrorExample() *customError.CustomError {
	return customError.NewCustomError(customError.WithInternalServerError("???", "generic service error"))
}

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockService) {
	mockCtrl := gomock.NewController(t)
	mockService := mocks.NewMockService(mockCtrl)
	return mockCtrl, mockService
}

func TestCreateGroup_Success(t *testing.T) {
	w, ctx := functions.PrepareCtx("POST")
	group := models.CreateMockGroup()
	functions.SetReqBody(ctx, group)

	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	mockServices.EXPECT().CreateGroup(group).Return(group, nil)

	handler := NewGroupHandler(mockServices)
	handler.CreateGroup(ctx)

	var response models.Group
	functions.GetRespBody(w, &response)

	assert.Equal(t, ctx.Writer.Status(), http.StatusCreated)
	assert.Equal(t, group.Id, response.Id)
	assert.Equal(t, group.Name, response.Name)
	assert.Equal(t, len(group.Participants), len(response.Participants))
}

func TestCreateGroup_ServiceError(t *testing.T) {
	_, ctx := functions.PrepareCtx("POST")
	group := models.CreateMockGroup()
	functions.SetReqBody(ctx, group)

	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	mockServices.EXPECT().CreateGroup(group).Return(nil, internalErrorExample())

	handler := NewGroupHandler(mockServices)
	handler.CreateGroup(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusInternalServerError)
}

func TestCreateGroup_EmptyNameError(t *testing.T) {
	_, ctx := functions.PrepareCtx("POST")
	invalidGroup := models.CreateMockGroup()
	invalidGroup.Name = ""
	functions.SetReqBody(ctx, invalidGroup)

	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	handler := NewGroupHandler(mockServices)
	handler.CreateGroup(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusBadRequest)
}

func TestGetGroup_Success(t *testing.T) {
	w, ctx := functions.PrepareCtx("GET")
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	expectedGroup := models.CreateMockGroup()
	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	mockServices.EXPECT().GetGroupByID("1").Return(expectedGroup, nil)

	handler := NewGroupHandler(mockServices)
	handler.GetGroup(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusOK)
	var response models.Group
	functions.GetRespBody(w, &response)
	assert.Equal(t, expectedGroup.Id, response.Id)
	assert.Equal(t, expectedGroup.Name, response.Name)
	assert.Equal(t, len(expectedGroup.Participants), len(response.Participants))
}

func TestGetGroup_NotFound(t *testing.T) {
	_, ctx := functions.PrepareCtx("GET")
	ctx.Params = []gin.Param{{Key: "id", Value: "999"}}

	err := customError.NewCustomError(customError.WithNotFound("Group not found", "Not found"))
	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	mockServices.EXPECT().GetGroupByID("999").Return(nil, err)

	handler := NewGroupHandler(mockServices)
	handler.GetGroup(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusNotFound)
}

func TestDeleteGroup_Success(t *testing.T) {
	_, ctx := functions.PrepareCtx("DELETE")
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	mockServices.EXPECT().DeleteGroup("1").Return(nil)

	handler := NewGroupHandler(mockServices)
	handler.DeleteGroup(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusNoContent)
}

func TestDeleteGroup_NotFound(t *testing.T) {
	_, ctx := functions.PrepareCtx("DELETE")
	ctx.Params = []gin.Param{{Key: "id", Value: "999"}}

	err := customError.NewCustomError(customError.WithNotFound("Group not found", "Not Found"))
	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	mockServices.EXPECT().DeleteGroup("999").Return(err)

	handler := NewGroupHandler(mockServices)
	handler.DeleteGroup(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusNotFound)
}

func TestGetMyMatch_EmptyUsername(t *testing.T) {
	_, ctx := functions.PrepareCtx("GET")
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}
	ctx.Request.URL.RawQuery = "username="

	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	handler := NewGroupHandler(mockServices)
	handler.GetMyMatch(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusBadRequest)
}

func TestGetMyMatch_EmptyGroup(t *testing.T) {
	_, ctx := functions.PrepareCtx("GET")
	ctx.Request.URL.RawQuery = "username=testuser"

	mockCtrl, mockServices := setupTest(t)
	defer mockCtrl.Finish()

	handler := NewGroupHandler(mockServices)
	handler.GetMyMatch(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusBadRequest)
}
