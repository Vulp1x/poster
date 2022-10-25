package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	authservice "github.com/inst-api/poster/gen/auth_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/sessions"
	"github.com/inst-api/poster/internal/store/users"
	"github.com/inst-api/poster/pkg/logger"
	"goa.design/goa/v3/security"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDContextKey = contextKey("UserID")

type userStore interface {
	FindByLogin(ctx context.Context, email string) (dbmodel.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, login, hashedPass string) error
}

// auth_service service example implementation.
// The example methods log the requests and return zero values.
type AuthServicesrvc struct {
	store       userStore
	securityCfg sessions.Configuration
}

// NewAuthService returns the auth_service service implementation.
func NewAuthService(userStore userStore, cfg sessions.Configuration) *AuthServicesrvc {
	return &AuthServicesrvc{securityCfg: cfg, store: userStore}
}

// BasicAuth implements the authorization logic for service "auth_service" for
// the "basic" security scheme.
func (s *AuthServicesrvc) BasicAuth(ctx context.Context, login, pass string, scheme *security.BasicScheme) (context.Context, error) {
	u, err := s.store.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			logger.Debugf(ctx, "No user with login: %s", login)
		} else {
			logger.Errorf(ctx, "failed to find user with login %s: internal error: %v", login, err)
		}

		return ctx, authservice.Unauthorized("")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pass)); err != nil {
		return ctx, authservice.Unauthorized("")
	}

	ctx = logger.WithFields(ctx, logger.Fields{"user_id": u.ID})
	ctx = AddUserIDToContext(ctx, u.ID)

	logger.Debugf(ctx, "Successfully checked credentials for login: %s", u.ID.String())

	return ctx, nil
}

// JWTAuth implements the authorization logic for service "auth_service" for
// the "jwt" security scheme.
func (s *AuthServicesrvc) JWTAuth(ctx context.Context, tokenString string, scheme *security.JWTScheme) (context.Context, error) {
	if tokenString == "" {
		return ctx, authservice.BadRequest("No token provided")
	}

	token, err := jwt.ParseWithClaims(tokenString, &sessions.SessionClaims{}, s.securityCfg.KeyFunc)
	if err != nil {
		logger.Infof(ctx, "Failed to parse token: %v", err)

		return ctx, authservice.Unauthorized("failed to parse token")
	}

	// logger.Debugf(ctx, "got sec scheme %+v", scheme)

	if claims, ok := token.Claims.(*sessions.SessionClaims); ok && token.Valid {
		ctx = logger.WithFields(ctx, logger.Fields{"user_id": claims.UserID})
		ctx = AddUserIDToContext(ctx, claims.UserID)
		logger.Debugf(ctx, "Successfully checked token")
		return ctx, nil
	}

	logger.Infof(ctx, "Token is valid: %t or Claims are SessionClaims: %#v", token.Valid, reflect.TypeOf(token.Claims))

	return ctx, fmt.Errorf("internal err")
}

// Signin Creates a valid JWT
func (s *AuthServicesrvc) Signin(ctx context.Context, p *authservice.SigninPayload) (*authservice.Creds, error) {
	logger.Info(ctx, "authService.signin")

	userID, err := UserIDFromContext(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to get userID from ctx: %v", err)

		return nil, authservice.Unauthorized("")
	}

	signedToken, tokenErr := sessions.
		NewSession(userID, s.securityCfg.TokenDuration).
		GenerateSignedToken(s.securityCfg.SigningKey)
	if tokenErr != nil {
		logger.Errorf(ctx, "failed to create new token")
		return nil, authservice.Unauthorized("")
	}

	return &authservice.Creds{JWT: signedToken}, nil
}

func (s *AuthServicesrvc) Profile(ctx context.Context, payload *authservice.ProfilePayload) (err error) {
	logger.Infof(ctx, "auth profile method")

	return nil
}

// AddUserIDToContext добавляет user_id в контекст
func AddUserIDToContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("no user id in ctx: %#v", ctx)
	}

	return userID, nil
}
