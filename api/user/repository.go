package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetSalt() int
	FindByCredential(ctx context.Context, user *User) error
}

type dbRepository struct {
	db          *pgxpool.Pool
	BCRYPT_SALT int
}

func NewRepository(db *pgxpool.Pool, BCRYPT_SALT int) Repository {
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
	row := d.db.QueryRow(ctx, stmt, user.Name, user.Password, user.Credential)
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
	row := d.db.QueryRow(ctx, stmt, user.Credential)
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	return nil

}