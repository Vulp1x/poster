package dbmodel

import "github.com/google/uuid"

func (b BotAccount) GetID() uuid.UUID {
	return b.ID
}

func (t TargetUser) GetID() uuid.UUID {
	return t.ID
}
