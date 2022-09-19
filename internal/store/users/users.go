package users

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/store"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// ErrLoginIsAlreadyUsed repository error.
var ErrLoginIsAlreadyUsed = errors.New("same email already registered")

// ErrUserNotFound repository error.
var ErrUserNotFound = errors.New("not found user")

// ErrUsersVehicleNotFound repository error.
var ErrUsersVehicleNotFound = errors.New("not found user's vehicle")

// ErrWorkShiftNotFound repository error.
var ErrWorkShiftNotFound = errors.New("not found user's work shift")

// NewStore creates new user store
func NewStore(opts ...store.Option) UserStore {
	var cfg store.Configuration
	for _, opt := range opts {
		opt(&cfg)
	}

	return UserStore{cfg: cfg}
}

// UserStore is a store, that provide logic for operating with users in database
type UserStore struct {
	cfg store.Configuration
}

func (us UserStore) GenerateNewPassword(ctx context.Context, id uuid.UUID) (string, error) {
	password := generatePassword()

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash for pass: %v", err)
	}

	tx, err := us.cfg.TxFunc(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v : %w", err, store.ErrTransactionFail)
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	err = q.UpdateUserPassword(ctx, dbmodel.UpdateUserPasswordParams{
		PasswordHash: string(hashPass),
		ID:           id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf("failed to update User password: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return password, nil
}

var firstPartsOfPasswords = []string{
	"жёлтый", "синий", "красный", "зелёный", "чёрный", "белый", "алый",
}

var secondPartsOfPasswords = []string{
	"лев", "тигр", "волк", "бизон", "орел", "ястреб", "коршун", "аист", "ягуар", "гепард", "леопард",
}

func generatePassword() string {
	pass := fmt.Sprintf("%s-%s",
		firstPartsOfPasswords[rand.Intn(len(firstPartsOfPasswords))],
		secondPartsOfPasswords[rand.Intn(len(secondPartsOfPasswords))],
	)

	return pass
}

// FindByID returns user by id
func (us UserStore) FindByID(ctx context.Context, id uuid.UUID) (UserProfile, error) {
	q := dbmodel.New(us.cfg.DBTXf(ctx))

	u, err := q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserProfile{}, ErrUserNotFound
		}

		return UserProfile{}, err
	}

	return UserProfile{User: u}, nil
}

// FindByLogin returns user with provided email.
// If there are now user with provided email, then ErrUserNotFound is returned.
func (us UserStore) FindByLogin(ctx context.Context, login string) (dbmodel.User, error) {
	q := dbmodel.New(us.cfg.DBTXf(ctx))

	u, err := q.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbmodel.User{}, ErrUserNotFound
		}

		return dbmodel.User{}, err
	}

	return u, nil
}

// Create ...
func (us UserStore) Create(ctx context.Context, login, passwordHash string) error {
	q := dbmodel.New(us.cfg.DBTXf(ctx))

	_, err := q.CreateUser(ctx, dbmodel.CreateUserParams{
		Login:        login,
		PasswordHash: passwordHash,
		Role:         0,
	})
	if err != nil {
		if strings.Contains(err.Error(), "login") {
			return ErrLoginIsAlreadyUsed
		}

		return err
	}

	return nil
}

func (us UserStore) Delete(ctx context.Context, id uuid.UUID) error {
	q := dbmodel.New(us.cfg.DBTXf(ctx))

	err := q.DeleteUserByID(ctx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}

	return err
}
