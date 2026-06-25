package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"stops/backend/internal/models"
	"stops/backend/internal/utils"
)

type UserRepository struct {
	db *pgxpool.Pool
}

type CreateUserParams struct {
	Name         string
	Email        *string
	Phone        *string
	PasswordHash string
	Role         models.UserRole
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) Create(ctx context.Context, params CreateUserParams) (models.User, error) {
	const query = `
		INSERT INTO users (name, email, phone, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, name, email, phone, password_hash, role, contribution_score, created_at, updated_at
	`

	user, err := scanUser(repository.db.QueryRow(ctx, query, params.Name, params.Email, params.Phone, params.PasswordHash, params.Role))
	if err != nil {
		return models.User{}, mapCreateUserError(err)
	}

	return user, nil
}

func (repository *UserRepository) FindByEmailOrPhone(ctx context.Context, email *string, phone *string) (models.User, error) {
	const query = `
		SELECT id::text, name, email, phone, password_hash, role, contribution_score, created_at, updated_at
		FROM users
		WHERE ($1::text IS NOT NULL AND email = $1)
		   OR ($2::text IS NOT NULL AND phone = $2)
		LIMIT 1
	`

	user, err := scanUser(repository.db.QueryRow(ctx, query, email, phone))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, utils.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("find user by identifier: %w", err)
	}

	return user, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	const query = `
		SELECT id::text, name, email, phone, password_hash, role, contribution_score, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	user, err := scanUser(repository.db.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, utils.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (models.User, error) {
	var user models.User
	var role string

	err := scanner.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&role,
		&user.ContributionScore,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	user.Role = models.UserRole(role)
	return user, nil
}

func mapCreateUserError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return utils.ErrEmailAlreadyInUse
		case "users_phone_key":
			return utils.ErrPhoneAlreadyInUse
		}
	}

	return fmt.Errorf("create user: %w", err)
}

func (repository *UserRepository) UpdateResetToken(ctx context.Context, userID string, token string, expiry time.Time) error {
	const query = `
		UPDATE users
		SET reset_token = $1, reset_token_expiry = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := repository.db.Exec(ctx, query, token, expiry, userID)
	if err != nil {
		return fmt.Errorf("update user reset token: %w", err)
	}
	return nil
}

func (repository *UserRepository) FindByResetToken(ctx context.Context, token string) (models.User, error) {
	const query = `
		SELECT id::text, name, email, phone, password_hash, role, contribution_score, created_at, updated_at
		FROM users
		WHERE reset_token = $1
		  AND reset_token_expiry > NOW()
		LIMIT 1
	`
	user, err := scanUser(repository.db.QueryRow(ctx, query, token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errors.New("reset token is invalid or expired")
		}
		return models.User{}, fmt.Errorf("find user by reset token: %w", err)
	}
	return user, nil
}

func (repository *UserRepository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	const query = `
		UPDATE users
		SET password_hash = $1, reset_token = NULL, reset_token_expiry = NULL, updated_at = NOW()
		WHERE id = $2
	`
	_, err := repository.db.Exec(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	return nil
}
