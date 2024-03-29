package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/andrii-stp/users-crud/model"
)

type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id int64, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type PostgresUserRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

var (
	_ UserRepository = (*PostgresUserRepository)(nil)

	ErrAlreadyExist = errors.New("username already in use")
	ErrUserNotFound = errors.New("user don't exist")
)

func NewPostgresRepository(logger *slog.Logger, db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		logger: logger,
		db:     db,
	}
}

func (ps PostgresUserRepository) List(ctx context.Context) ([]model.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rows, err := psql.Select("*").From("users").RunWith(ps.db).QueryContext(ctx)
	if err != nil {
		ps.logger.Error("Failed to execute select query",
			slog.String("err", err.Error()))

		return nil, err
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName,
			&user.Email, &user.Status, &user.Department); err != nil {
			return nil, fmt.Errorf("failed to scan select data: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ps PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	exist, err := getByID(ctx, tx, user.UserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if exist != nil {
		return ErrAlreadyExist
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	err = psql.Insert("users").
		Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
		Values(user.UserName, user.FirstName, user.LastName, user.Email, user.Status, user.Department).
		Suffix("RETURNING *").RunWith(tx).QueryRowContext(ctx).
		Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.Department)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ps PostgresUserRepository) Update(ctx context.Context, id int64, user *model.User) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("update user transaction failed: %w", err)
	}

	defer tx.Rollback()

	// get existing user by id
	var targeted *model.User

	targeted, err = getByID(ctx, tx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if targeted == nil {
		return ErrUserNotFound
	}

	// check requested username if it's alredy in use
	exist, err := getByUserName(ctx, tx, user.UserName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if exist != nil && exist.UserID != targeted.UserID {
		return ErrAlreadyExist
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	err = psql.Update("users").SetMap(
		sq.Eq{
			"user_name":   user.UserName,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"user_status": user.Status,
			"department":  user.Department,
		}).Suffix("RETURNING *").RunWith(tx).QueryRowContext(ctx).
		Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.Department)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ps PostgresUserRepository) Delete(ctx context.Context, id int64) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// get by username targeted user
	var exist *model.User

	exist, err = getByID(ctx, tx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if exist == nil {
		return ErrUserNotFound
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	_, err = psql.Delete("users").Where(sq.Eq{"user_id": id}).RunWith(tx).ExecContext(ctx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func getByID(ctx context.Context, tx *sql.Tx, id int64) (*model.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var user model.User
	err := psql.Select("*").From("users").Where(sq.Eq{"user_id": id}).
		RunWith(tx).QueryRowContext(ctx).
		Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName,
			&user.Email, &user.Status, &user.Department)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getByUserName(ctx context.Context, tx *sql.Tx, username string) (*model.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var user model.User
	err := psql.Select("*").From("users").Where(sq.Eq{"user_name": username}).
		RunWith(tx).QueryRowContext(ctx).
		Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName,
			&user.Email, &user.Status, &user.Department)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
