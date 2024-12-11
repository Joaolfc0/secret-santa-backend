package models

import (
	"time"

	"github.com/invopop/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty" swaggerignore:"true"`
	Name         string             `json:"name" bson:"name" example:"Equipe pe no chao"`
	Participants []Participant      `json:"participants" bson:"participants" example:"joao,mari"`
	Matches      []Match            `json:"matches" bson:"matches"  swaggerignore:"true"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt, omitempty" swaggerignore:"true"`
	UpdatedAt    time.Time          `json:"updateAt" bson:"updateAt, omitempty" swaggerignore:"true"`
}

type Participant struct {
	Name  string `json:"name" bson:"name" example:"Mari"`
	Email string `json:"email" bson:"email" example:"Mari@gmail.com"`
}

type Match struct {
	First  string `json:"first" bson:"first"  example:"joao"`
	Second string `json:"second" bson:"second" example:"mari"`
}

func (l Group) Validate() error {
	err := validation.ValidateStruct(&l,
		validation.Field(&l.Name, validation.Required),
	)

	if err != nil {
		return err
	}

	return nil
}

func (l Participant) Validate() error {
	err := validation.ValidateStruct(&l,
		validation.Field(&l.Name, validation.Required),
		validation.Field(&l.Email, validation.Required),
	)

	if err != nil {
		return err
	}

	return nil
}
