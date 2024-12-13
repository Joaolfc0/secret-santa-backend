package group

import (
	"testing"

	"service-secret-santa/config"
	"service-secret-santa/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func groupToBSON(g *models.Group) bson.D {
	temp, _ := bson.Marshal(g)
	var groupBSON bson.D
	_ = bson.Unmarshal(temp, &groupBSON)
	return groupBSON
}

func TestGetMyMatch(t *testing.T) {
	config.LoadConfig()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewGroupRepository(mt.Client)

		group := models.Group{
			Id: primitive.NewObjectID(),
			Matches: []models.Match{
				{First: "João", Second: "Mario"},
				{First: "Mario", Second: "Luigi"},
				{First: "Luigi", Second: "João"},
			},
		}
		groupBSON := groupToBSON(&group)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "secret-santa.groups", mtest.FirstBatch, groupBSON),
			mtest.CreateCursorResponse(0, "secret-santa.groups", mtest.FirstBatch, groupBSON),
			mtest.CreateCursorResponse(0, "secret-santa.groups", mtest.FirstBatch, groupBSON),
		)

		match, err := repo.GetMyMatch(group.Id.Hex(), "João")
		assert.Nil(t, err)
		assert.Equal(t, match, "Mario")

		match, err = repo.GetMyMatch(group.Id.Hex(), "Mario")
		assert.Nil(t, err)
		assert.Equal(t, match, "Luigi")

		match, err = repo.GetMyMatch(group.Id.Hex(), "Luigi")
		assert.Nil(t, err)
		assert.Equal(t, match, "João")
	})

	mt.Run("no match", func(mt *mtest.T) {
		repo := NewGroupRepository(mt.Client)

		group := models.Group{
			Id: primitive.NewObjectID(),
			Matches: []models.Match{
				{First: "João", Second: "Mario"},
			},
		}
		groupBSON := groupToBSON(&group)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "secret-santa.groups", mtest.FirstBatch, groupBSON),
		)

		_, err := repo.GetMyMatch(group.Id.Hex(), "Mario")
		assert.Equal(t, err.Status, 404)
	})

	mt.Run("db error", func(mt *mtest.T) {
		repo := NewGroupRepository(mt.Client)

		mt.AddMockResponses(
			mtest.CreateWriteErrorsResponse(mtest.WriteError{Code: 11000}),
		)

		_, err := repo.GetMyMatch(primitive.NewObjectID().Hex(), "João")
		assert.Equal(t, err.Status, 500)
	})
}
