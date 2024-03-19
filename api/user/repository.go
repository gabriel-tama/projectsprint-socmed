package user

import (
	"context"
	"errors"

	"github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetSalt() int
	FindByCredential(ctx context.Context, user *User) error
	AddEmail(ctx context.Context, email string, user_id int) error
	AddPhone(ctx context.Context, phone string, user_id int) error
	UpdateAccount(ctx context.Context, name string, emailUrl string, user_id int) error
}

type dbRepository struct {
	db          *db.DB
	BCRYPT_SALT int
}

func NewRepository(db *db.DB, BCRYPT_SALT int) Repository {
	return &dbRepository{db: db, BCRYPT_SALT: BCRYPT_SALT}
}

func (d *dbRepository) GetSalt() int {
	return d.BCRYPT_SALT
}

func (d *dbRepository) Create(ctx context.Context, user *User) error {

	stmt := `
        INSERT INTO users (
            name, password, 
    `
	switch user.CredentialType {
	case "email":
		stmt = stmt + "email) "
	case "phone":
		stmt = stmt + "phoneNumber) "
	}
	stmt = stmt + `VALUES ($1, $2, $3) RETURNING id`
	row := d.db.Pool.QueryRow(ctx, stmt, user.Name, user.Password, user.Credential)
	err := row.Scan(&user.ID)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrUsernameAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (d *dbRepository) FindByCredential(ctx context.Context, user *User) error {
	stmt := `SELECT id, name, password `
	switch user.CredentialType {
	case "email":
		stmt = stmt + "email " + "FROM users WHERE email"
	case "phone":
		stmt = stmt + "phoneNumber " + "FROM users WHERE phoneNumber"
	}
	stmt = stmt + "=$1 "
	row := d.db.Pool.QueryRow(ctx, stmt, user.Credential)
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	return nil

}

func (d *dbRepository) AddEmail(ctx context.Context, email string, user_id int) error {
	result, err := d.db.Pool.Exec(ctx, "UPDATE users SET email=$1 WHERE id=$2 AND email IS NULL", email, user_id)
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrEmailAlreadyExists
	}

	if result.RowsAffected() == 0 {
		return ErrWrongRoute
	}

	return err

}

func (d *dbRepository) AddPhone(ctx context.Context, phone string, user_id int) error {
	result, err := d.db.Pool.Exec(ctx, "UPDATE users SET phoneNumber=$1 WHERE id=$2 AND phoneNumber IS NULL", phone, user_id)
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrPhoneAlreadyExists
	}

	if result.RowsAffected() == 0 {
		return ErrWrongRoute
	}

	return err

}

func (d *dbRepository) UpdateAccount(ctx context.Context, name string, imageUrl string, user_id int) error {
	_, err := d.db.Pool.Exec(ctx, "UPDATE users SET name=$1, imageUrl=$2 WHERE id=$3", name, imageUrl, user_id)

	return err

}
