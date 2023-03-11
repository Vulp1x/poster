package mw

import (
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/inst-api/poster/pkg/logger"
)

const (
	// bodyRequestKey key for use in context.
	bodyRequestKey contextKey = "Body"
	// idRequestKey key for use in context.
	idRequestKey contextKey = "ID Query"
)

// ReadBody read body from http middleware.
func ReadBody() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logger.Errorf(r.Context(), "Failed to read body: %s", err)
				http.Error(w, "This request needs body.", http.StatusBadRequest)

				return
			}

			defer r.Body.Close()

			logger.Debugf(r.Context(), "Body successfully read")
			next.ServeHTTP(w, bodyRequestKey.Write(r, body))
		}

		return http.HandlerFunc(fn)
	}
}

// ReadIDPath reads id from url path.
func ReadIDPath() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				BadRequest(w, r, "failed to parse id: %v", err)

				return
			}

			logger.Debugf(r.Context(), "ID successfully read %v", id)

			next.ServeHTTP(w, idRequestKey.Write(r, id))
		}

		return http.HandlerFunc(fn)
	}
}

// GetBody is used to get request's body from context.
func GetBody(r *http.Request) []byte {
	body, ok := r.Context().Value(bodyRequestKey).([]byte)
	if !ok {
		panic("can`t get body from context when expected" +
			string(debug.Stack()))
	}

	return body
}

// GetIDFromPath is used to get id in path from context.
func GetIDFromPath(r *http.Request) uuid.UUID {
	id, ok := r.Context().Value(idRequestKey).(uuid.UUID)
	if !ok {
		panic("can`t get id query from context when expected")
	}

	return id
}
