package instagrapi

import (
	"math/rand"
	"strconv"

	"github.com/inst-api/poster/internal/dbmodel"
)

type UserShort struct {
	Pk       string `json:"pk"`
	Username string `json:"username"`
}

type UserTag struct {
	User UserShort `json:"user"`
	X    float64   `json:"x"`
	Y    float64   `json:"y"`
}

func prepareUserTags(users []dbmodel.TargetUser) []UserTag {
	tags := make([]UserTag, len(users))
	for i, user := range users {
		tags[i] = UserTag{
			User: UserShort{Pk: strconv.FormatInt(user.UserID, 10), Username: user.Username},
			X:    rand.Float64(),
			Y:    rand.Float64(),
		}
	}

	return tags
}
