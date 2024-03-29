// Code generated by goa v3.11.3, DO NOT EDIT.
//
// auth_service endpoints
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package authservice

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "auth_service" service endpoints.
type Endpoints struct {
	Signin  goa.Endpoint
	Profile goa.Endpoint
}

// NewEndpoints wraps the methods of the "auth_service" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Signin:  NewSigninEndpoint(s, a.BasicAuth),
		Profile: NewProfileEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "auth_service" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Signin = m(e.Signin)
	e.Profile = m(e.Profile)
}

// NewSigninEndpoint returns an endpoint function that calls the method
// "signin" of service "auth_service".
func NewSigninEndpoint(s Service, authBasicFn security.AuthBasicFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*SigninPayload)
		var err error
		sc := security.BasicScheme{
			Name:           "basic",
			Scopes:         []string{"driver"},
			RequiredScopes: []string{},
		}
		ctx, err = authBasicFn(ctx, p.Login, p.Password, &sc)
		if err != nil {
			return nil, err
		}
		return s.Signin(ctx, p)
	}
}

// NewProfileEndpoint returns an endpoint function that calls the method
// "profile" of service "auth_service".
func NewProfileEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*ProfilePayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.Profile(ctx, p)
	}
}
