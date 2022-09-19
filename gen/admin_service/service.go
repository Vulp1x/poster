// Code generated by goa v3.8.5, DO NOT EDIT.
//
// admin_service service
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package adminservice

import (
	"context"

	"goa.design/goa/v3/security"
)

// The secured service exposes endpoints that require valid authorization
// credentials.
type Service interface {
	// admins could add drivers from main system
	AddManager(context.Context, *AddManagerPayload) (err error)
	// admins could delete managers from main system
	DropManager(context.Context, *DropManagerPayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "admin_service"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"add_manager", "drop_manager"}

// AddManagerPayload is the payload type of the admin_service service
// add_manager method.
type AddManagerPayload struct {
	// JWT used for authentication
	Token    *string
	Login    string
	Password string
}

// DropManagerPayload is the payload type of the admin_service service
// drop_manager method.
type DropManagerPayload struct {
	// JWT used for authentication
	Token *string
	// id менеджера, которого необходимо удалить
	ManagerID string `json:"manager_id"`
}

// Invalid request
type BadRequest string

// internal error
type InternalError string

// Token scopes are invalid
type InvalidScopes string

// Credentials are invalid
type Unauthorized string

type UserNotFound string

// Error returns an error description.
func (e BadRequest) Error() string {
	return "Invalid request"
}

// ErrorName returns "bad request".
func (e BadRequest) ErrorName() string {
	return "bad request"
}

// Error returns an error description.
func (e InternalError) Error() string {
	return "internal error"
}

// ErrorName returns "internal error".
func (e InternalError) ErrorName() string {
	return "internal error"
}

// Error returns an error description.
func (e InvalidScopes) Error() string {
	return "Token scopes are invalid"
}

// ErrorName returns "invalid-scopes".
func (e InvalidScopes) ErrorName() string {
	return "invalid-scopes"
}

// Error returns an error description.
func (e Unauthorized) Error() string {
	return "Credentials are invalid"
}

// ErrorName returns "unauthorized".
func (e Unauthorized) ErrorName() string {
	return "unauthorized"
}

// Error returns an error description.
func (e UserNotFound) Error() string {
	return ""
}

// ErrorName returns "user not found".
func (e UserNotFound) ErrorName() string {
	return "user not found"
}
