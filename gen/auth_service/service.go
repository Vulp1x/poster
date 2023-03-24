// Code generated by goa v3.11.3, DO NOT EDIT.
//
// auth_service service
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package authservice

import (
	"context"

	"goa.design/goa/v3/security"
)

// The secured service exposes endpoints that require valid authorization
// credentials.
type Service interface {
	// Creates a valid JWT
	Signin(context.Context, *SigninPayload) (res *Creds, err error)
	// get user profile
	Profile(context.Context, *ProfilePayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// BasicAuth implements the authorization logic for the Basic security scheme.
	BasicAuth(ctx context.Context, user, pass string, schema *security.BasicScheme) (context.Context, error)
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "auth_service"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"signin", "profile"}

// Creds is the result type of the auth_service service signin method.
type Creds struct {
	// JWT token
	JWT string
}

// ProfilePayload is the payload type of the auth_service service profile
// method.
type ProfilePayload struct {
	// JWT used for authentication
	Token string
}

// Credentials used to authenticate to retrieve JWT token
type SigninPayload struct {
	// login used to perform signin
	Login string
	// Password used to perform signin
	Password string
}

// Invalid request
type BadRequest string

// internal error
type InternalError string

// Credentials are invalid
type Unauthorized string

// Not found
type UserNotFound string

// Error returns an error description.
func (e BadRequest) Error() string {
	return "Invalid request"
}

// ErrorName returns "bad request".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e BadRequest) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "bad request".
func (e BadRequest) GoaErrorName() string {
	return "bad request"
}

// Error returns an error description.
func (e InternalError) Error() string {
	return "internal error"
}

// ErrorName returns "internal error".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e InternalError) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "internal error".
func (e InternalError) GoaErrorName() string {
	return "internal error"
}

// Error returns an error description.
func (e Unauthorized) Error() string {
	return "Credentials are invalid"
}

// ErrorName returns "unauthorized".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e Unauthorized) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "unauthorized".
func (e Unauthorized) GoaErrorName() string {
	return "unauthorized"
}

// Error returns an error description.
func (e UserNotFound) Error() string {
	return "Not found"
}

// ErrorName returns "user not found".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e UserNotFound) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "user not found".
func (e UserNotFound) GoaErrorName() string {
	return "user not found"
}
