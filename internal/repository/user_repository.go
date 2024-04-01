package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/paimon_bank/internal/entity"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(request *entity.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`

	_, err := r.DB.Exec(query, request.Name, request.Email, request.Password)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	user := new(entity.User)
	query := `SELECT id, name, email, password from users where email = $1`

	err := r.DB.Get(user, query, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
