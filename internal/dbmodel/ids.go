package dbmodel

import "github.com/google/uuid"

func (b BotAccount) GetID() uuid.UUID {
	return b.ID
}
