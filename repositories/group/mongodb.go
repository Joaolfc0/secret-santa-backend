package group

import (
	"context"
	"service-secret-santa/config"
	"service-secret-santa/customError"
	"service-secret-santa/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateGroup(group *models.Group) (*models.Group, *customError.CustomError)
	GetGroupByID(id string) (*models.Group, *customError.CustomError)
	UpdateGroup(id string, group *models.Group) (*models.Group, *customError.CustomError)
	DeleteGroup(id string) *customError.CustomError
	AddParticipant(id string, participant *models.Participant) (*models.Group, *customError.CustomError)
	UpdateMatches(id string, matches []models.Match) *customError.CustomError
	GetAllGroups() ([]*models.Group, *customError.CustomError)
	GetMyMatch(id string, username string) (string, *customError.CustomError)
}

type resource struct {
	db *mongo.Client
}

func NewGroupRepository(db *mongo.Client) Repository {
	return &resource{db: db}
}

func (r *resource) CreateGroup(group *models.Group) (*models.Group, *customError.CustomError) {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")
	result, err := collection.InsertOne(context.Background(), group)
	if err != nil {
		return nil, customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Failed to create group"))
	}

	group.Id = result.InsertedID.(primitive.ObjectID)
	return group, nil
}

func (r *resource) GetGroupByID(id string) (*models.Group, *customError.CustomError) {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, customError.NewCustomError(customError.WithBadRequest("Invalid group ID", "Invalid ID format"))
	}

	var group models.Group
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, customError.NewCustomError(customError.WithNotFound("Group not found", "No group found with the given ID"))
		}
		return nil, customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Error finding group"))
	}

	return &group, nil
}

func (r *resource) UpdateGroup(id string, group *models.Group) (*models.Group, *customError.CustomError) {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, customError.NewCustomError(customError.WithBadRequest("Invalid group ID", "Invalid ID format"))
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": group})
	if err != nil {
		return nil, customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Failed to update group"))
	}

	return group, nil
}

func (r *resource) DeleteGroup(id string) *customError.CustomError {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customError.NewCustomError(customError.WithBadRequest("Invalid group ID", "Invalid ID format"))
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Failed to delete group"))
	}

	return nil
}

func (r *resource) AddParticipant(id string, participant *models.Participant) (*models.Group, *customError.CustomError) {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, customError.NewCustomError(customError.WithBadRequest("Invalid group ID", "Invalid ID format"))
	}

	update := bson.M{"$addToSet": bson.M{"participants": participant}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Failed to add participant"))
	}

	return r.GetGroupByID(id)
}

func (r *resource) UpdateMatches(id string, matches []models.Match) *customError.CustomError {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customError.NewCustomError(customError.WithBadRequest("Invalid group ID", "Invalid ID format"))
	}

	update := bson.M{"$set": bson.M{"matches": matches}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Failed to update matches"))
	}

	return nil
}

func (r *resource) GetMyMatch(id string, username string) (string, *customError.CustomError) {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", customError.NewCustomError(customError.WithBadRequest("Invalid group ID", "Invalid ID format"))
	}

	var group models.Group
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", customError.NewCustomError(customError.WithNotFound("Group not found", "No group found with the given ID"))
		}
		return "", customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Error finding group"))
	}

	for _, match := range group.Matches {
		if match.First == username {
			return match.Second, nil
		}
	}

	return "", customError.NewCustomError(customError.WithNotFound("Match not found", "No match found for the given username"))
}

func (r *resource) GetAllGroups() ([]*models.Group, *customError.CustomError) {
	collection := r.db.Database(config.Cfg.MongoDB).Collection("groups")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Error retrieving groups"))
	}
	defer cursor.Close(context.Background())

	var groups []*models.Group
	if err = cursor.All(context.Background(), &groups); err != nil {
		return nil, customError.NewCustomError(customError.WithInternalServerError(err.Error(), "Error decoding groups"))
	}

	return groups, nil
}
