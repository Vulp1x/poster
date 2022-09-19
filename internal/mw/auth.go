package mw

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"runtime/debug"

	"github.com/dgrijalva/jwt-go"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/sessions"
	"github.com/inst-api/poster/pkg/logger"
)

type contextKey string

func (c contextKey) String() string {
	return "middlewares context value " + string(c)
}

func (c contextKey) Write(r *http.Request, val interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), c, val))
}

const (
	// claimsRequestKey key for using in context.
	claimsRequestKey contextKey = "Claims"
	// tokenRequestKey key for using in context.
	tokenRequestKey contextKey = "Token"
)

// CheckSession check sesssion middleware.
func CheckSession(securityConfig sessions.Configuration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Bearer")
			if tokenString == "" {
				Error(w, r, http.StatusUnauthorized, "No token provided")

				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &sessions.SessionClaims{}, securityConfig.KeyFunc)
			if err != nil {
				InternalError(w, r, "Failed to parse token: %v", err)

				return
			}

			if claims, ok := token.Claims.(*sessions.SessionClaims); ok && token.Valid {
				LogEntrySetField(r, "user_id", claims.UserID)
				logger.Debugf(r.Context(), "Successfully checked token")

				next.ServeHTTP(w, claimsRequestKey.Write(r, claims))
			} else {
				InternalError(w, r, "Token is valid: %t or Claims are SessionClaims: %v", token.Valid, reflect.TypeOf(token.Claims))

				return
			}
		}

		return http.HandlerFunc(fn)
	}
}

// CheckCredentials check user credentials middleware.
func CheckCredentials(txFunc dbmodel.DBTXFunc, conf sessions.Configuration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			body := GetBody(r)

			bodyUser := make(map[string]string)

			err := json.Unmarshal(body, &bodyUser)
			if err != nil {
				BadRequest(w, r, "Failed to unmarshal User json: %s ", err)

				return
			}

			ctx := r.Context()
			q := dbmodel.New(txFunc(ctx))

			u, err2 := q.FindByEmail(ctx, bodyUser["email"])
			if err2 != nil {
				logger.Infof(ctx, "No user with email: %s", bodyUser["email"])
				http.Error(w, "No user with email: "+bodyUser["email"], http.StatusNotFound)

				return
			}

			if bodyUser["password"] != u.PasswordHash {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			LogEntrySetField(r, "user_id", u.ID)

			signedToken, tokenErr := sessions.
				NewSession(u.ID, conf.TokenDuration).
				GenerateSignedToken(conf.SigningKey)
			if tokenErr != nil {
				InternalError(w, r, "Failed to create signed token: %v", tokenErr)

				return
			}

			logger.Debugf(ctx, "Successfully checked credentials for user: %s", u.ID.String())

			next.ServeHTTP(w, tokenRequestKey.Write(r, signedToken))
		}

		return http.HandlerFunc(fn)
	}
}

// GetToken is used to get User JSON Web Token from request context.
func GetToken(r *http.Request) string {
	signedToken, ok := r.Context().Value(tokenRequestKey).(string)
	if !ok {
		panic("can`t get signed token from context when expected" +
			string(debug.Stack()))
	}

	return signedToken
}

// GetClaims is used to get claims(userID etc) from request context.
func GetClaims(r *http.Request) *sessions.SessionClaims {
	claims, ok := r.Context().Value(claimsRequestKey).(*sessions.SessionClaims)
	if !ok {
		panic("can`t get claims from context when expected" +
			string(debug.Stack()))
	}

	return claims
}
