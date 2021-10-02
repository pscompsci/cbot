package pg

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/pscompsci/cbot/internal/repository"
	"github.com/pscompsci/cbot/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, password_hash, created_at, activated)
		VALUES($1, $2, $3, current_timestamp, TRUE)`

	_, err = m.DB.Exec(stmt, name, email, string(hash))
	if err != nil {
		var sqlError *pq.Error
		if errors.As(err, &sqlError) {
			if strings.Contains(sqlError.Message, "unique_email") {
				return repository.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hash []byte
	var activated bool

	stmt := `SELECT id, password_hash, activated FROM users WHERE email = $1`
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hash, &activated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repository.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, repository.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if !activated {
		return 0, repository.ErrUserNotActivated
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT * FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}
