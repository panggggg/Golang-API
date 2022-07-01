package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	BlogID				primitive.ObjectID			`json:"id,omitempty"`
	Title				string						`json:"title,omitempty" validate:"required"`
	Details				string						`json:"details,omitempty" validate:"required"`
	Author				string						`json:"author,omitempty" validate:"required"`
	Created_at			time.Time					`json:"created_at,omitempty"`
}