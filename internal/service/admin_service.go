package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	adminservice "github.com/inst-api/poster/gen/admin_service"
	authservice "github.com/inst-api/poster/gen/auth_service"
	"github.com/inst-api/poster/internal/store/users"
	"github.com/inst-api/poster/pkg/logger"
	"goa.design/goa/v3/security"
	"golang.org/x/crypto/bcrypt"
)

// admin_service service example implementation.
// The example methods log the requests and return zero values.
type adminServicesrvc struct {
	auth      authservice.Auther
	userStore userStore
}

// NewAdminService returns the admin_service service implementation.
func NewAdminService(auth authservice.Auther, store userStore) adminservice.Service {
	return &adminServicesrvc{auth: auth, userStore: store}
}

func (s *adminServicesrvc) DropManager(ctx context.Context, payload *adminservice.DropManagerPayload) error {
	managerID, err := uuid.Parse(payload.ManagerID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse id from '%s': %v", payload.ManagerID, err)
		return adminservice.BadRequest("failed to parse manager_id")
	}

	err = s.userStore.Delete(ctx, managerID)
	if err != nil {
		logger.Errorf(ctx, "failed to delete user with id '%s': %v", managerID, err)
		if errors.Is(err, users.ErrUserNotFound) {
			return adminservice.UserNotFound("")
		}

		return adminservice.InternalError("")
	}

	return nil
}

// JWTAuth implements the authorization logic for service "admin_service" for
// the "jwt" security scheme.
func (s *adminServicesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	//
	// TBD: add authorization logic.
	//
	// In case of authorization failure this function should return
	// one of the generated error structs, e.g.:
	//
	//    return ctx, myservice.MakeUnauthorizedError("invalid token")
	//
	// Alternatively this function may return an instance of
	// goa.ServiceError with a Name field value that matches one of
	// the design error names, e.g:
	//
	//    return ctx, goa.PermanentError("unauthorized", "invalid token")
	//

	return s.auth.JWTAuth(ctx, token, scheme)
}

// AddManager Add Manager by unique_id. Only admins could add Managers
func (s *adminServicesrvc) AddManager(ctx context.Context, p *adminservice.AddManagerPayload) error {
	logger.Debug(ctx, "admin_service.AddManager")

	hashPass, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(ctx, "failed to generate hash for pass: %v", err)
		return adminservice.InternalError("")
	}

	err = s.userStore.Create(ctx, p.Login, string(hashPass))
	if err != nil {
		logger.Errorf(ctx, "failed to create user: %v", err)
		return adminservice.InternalError("")
	}

	return nil
}
