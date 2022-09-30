package domain

import (
	"fmt"

	"github.com/google/uuid"
)

func GetIDs[T IDer](models []T) []uuid.UUID {
	ids := make([]uuid.UUID, len(models))
	for i, model := range models {
		ids[i] = model.GetID()
	}

	return ids
}

func Strings[T fmt.Stringer](models []T) []string {
	stringsToReturn := make([]string, len(models))
	for i, model := range models {
		stringsToReturn[i] = model.String()
	}

	return stringsToReturn
}
