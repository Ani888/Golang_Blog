package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type Article struct {
    Id       primitive.ObjectID     `json:"id,omitempty"`
    Title       string             `json:"title,omitempty" validate:"required"`
    Content     string             `json:"content,omitempty" validate:"required"`
    
}