package group

import (
	"strconv"
	"testing"
	"time"

	"service-secret-santa/customError"
	"service-secret-santa/models"
	mocks "service-secret-santa/repositories/group/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func internalErrorExample() *customError.CustomError {
	return customError.NewCustomError(customError.WithInternalServerError("???", "generic service error"))
}

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockRepository) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	return mockCtrl, mockRepo
}

func MockUnmatchedGroup(n int) *models.Group {
	g := &models.Group{
		Id:           primitive.NewObjectID(),
		Name:         "Test Group",
		Participants: []models.Participant{},
		CreatedAt:    time.Date(2023, 12, 11, 0, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2023, 12, 11, 0, 0, 0, 0, time.UTC),
	}
	for i := 0; i < n; i++ {
		i_str := strconv.Itoa(i)
		g.Participants = append(g.Participants, models.Participant{
			Name:  "P" + i_str,
			Email: "p" + i_str + "@gmail.com",
		})
	}
	return g
}

func TestMatchParticipants_Success(t *testing.T) {
	for i := 2; i < 70; i++ {
		mockCtrl, mockRepo := setupTest(t)
		service := NewGroupService(mockRepo)

		group := MockUnmatchedGroup(i)

		mockRepo.EXPECT().GetGroupByID(group.Id.Hex()).Return(group, nil)
		mockRepo.EXPECT().UpdateMatches(group.Id.Hex(), gomock.Any()).Return(nil)

		_, err := service.MatchParticipants(group.Id.Hex())

		assert.Nil(t, err)
		assert.Equal(t, len(group.Matches), i)

		for _, m := range group.Matches {
			assert.True(t, m.First != m.Second)
		}

		mockCtrl.Finish()
	}
}

func TestMatchParticipants_NotEnoughParticipants(t *testing.T) {
	mockCtrl, mockRepo := setupTest(t)
	defer mockCtrl.Finish()
	service := NewGroupService(mockRepo)

	group := MockUnmatchedGroup(1)

	mockRepo.EXPECT().GetGroupByID(group.Id.Hex()).Return(group, nil)

	_, err := service.MatchParticipants(group.Id.Hex())

	assert.Equal(t, err.Status, 400)
}

func TestMatchParticipants_DBError(t *testing.T) {
	mockCtrl, mockRepo := setupTest(t)
	defer mockCtrl.Finish()
	service := NewGroupService(mockRepo)

	group := MockUnmatchedGroup(2)
	mockErr := internalErrorExample()

	mockRepo.EXPECT().GetGroupByID(group.Id.Hex()).Return(nil, mockErr)

	_, err := service.MatchParticipants(group.Id.Hex())

	assert.Equal(t, err.Status, 500)
}
