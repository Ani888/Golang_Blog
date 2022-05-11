package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Comments  string             `json:"comment,omitempty" validate:"required"`
	ArticleId primitive.ObjectID `json:"articleId,omitempty" validate:"required"`
}
