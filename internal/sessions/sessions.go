package sessions

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Configuration represents token security.
type Configuration struct {
	SigningKey    string        `yaml:"JWT signing key"`
	TokenDuration time.Duration `yaml:"token validness duration"`
}

// Default sets default values in config variables.
func (c *Configuration) Default() {
	var hoursInYear time.Duration = 8760

	c.SigningKey = "SuperPuperKey42"
	c.TokenDuration = hoursInYear * time.Hour
}

var (
	// ErrSessionExpired auth error.
	ErrSessionExpired = errors.New("session is expired")
	// ErrUnexpectedSigningMethod auth error.
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

// SessionClaims session parameters to validate.
type SessionClaims struct {
	UserID     uuid.UUID `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	ValidUntil time.Time `json:"valid_until"`
}

// Valid validates session.
func (s *SessionClaims) Valid() error {
	if s.ValidUntil.Before(time.Now()) {
		return ErrSessionExpired
	}

	return nil
}

// NewSession creates new session.
func NewSession(userID uuid.UUID, duration time.Duration) *SessionClaims {
	return &SessionClaims{
		UserID:     userID,
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().Add(duration),
	}
}

func (s SessionClaims) String() string {
	return fmt.Sprintf("UserID: %d \nCreatedAt: %s, ValidUntil: %s",
		s.UserID, s.CreatedAt.Format(time.RFC1123), s.ValidUntil.Format(time.RFC1123))
}

// GenerateSignedToken generates jwt signed tokens.
func (s *SessionClaims) GenerateSignedToken(signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, s)

	return token.SignedString([]byte(signingKey))
}

// KeyFunc validates token after it was parsed.
func (c Configuration) KeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
	}

	return []byte(c.SigningKey), nil
}
