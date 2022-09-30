package domain

import "github.com/google/uuid"

// IDer is an interface for getting id of a model
type IDer interface {
	GetID() uuid.UUID
}
