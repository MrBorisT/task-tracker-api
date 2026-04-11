package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	Pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{Pool: pool}
}

func (s *UserStore) RegisterUser(ctx context.Context, userRequest models.RegisterUserRequest) error {
	trimmedEmail := strings.TrimSpace(userRequest.Email)
	if trimmedEmail == "" {
		return ErrEmptyUserEmail
	}

	trimmedPassword := strings.TrimSpace(userRequest.Password)
	if trimmedPassword == "" {
		return ErrEmptyUserPassword
	}
	if len(trimmedPassword) < 6 {
		return ErrShortUserPassword
	}
	if len(trimmedPassword) > 72 {
		return ErrLongUserPassword
	}

	hashedPassword, err := s.hashPassword(trimmedPassword)
	if err != nil {
		return err
	}

	newUser := models.User{
		ID:           s.generateID(),
		Email:        userRequest.Email,
		PasswordHash: hashedPassword,
	}
	query := "INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)"
	_, err = s.Pool.Exec(ctx, query, newUser.ID, newUser.Email, newUser.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PGCodeUniqueViolation {
			return ErrUserAlreadyExists
		}
	}
	return fmt.Errorf("create user: %w", err)
}

func (s *UserStore) generateID() string {
	return uuid.New().String()
}

func (s *UserStore) hashPassword(password string) (string, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptHash), nil
}
