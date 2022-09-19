package users

import "github.com/inst-api/poster/internal/dbmodel"

// UserProfile ...
type UserProfile struct {
	dbmodel.User
	// OnRoute показывает есть ли у водителя активная смена
	OnRoute bool
}
