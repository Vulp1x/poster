package domain

import (
	"fmt"
	"math/rand"

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

func RandomFromSlice[T interface{}](slice []T) T {
	switch len(slice) {
	case 0:
		panic(fmt.Sprintf("got empty slice %T", slice))
	case 1:
		return slice[0]
	default:
		return slice[rand.Intn(len(slice)-1)]
	}
}
