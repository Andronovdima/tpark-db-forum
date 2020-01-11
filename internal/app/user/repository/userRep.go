package repository

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(thisDB *sql.DB) *UserRepository {
	userRep := &UserRepository{
		db: thisDB,
	}
	return userRep
}

func (r *UserRepository) Create(user *models.User) error {
	_ , err := r.db.Exec(
		"INSERT INTO users about, email, fullname, nickname " +
			"VALUES $1, $2, $3, $4",
			user.About,
			user.Email,
			user.Fullname,
			user.Nickname ,
	)

	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return rerr
	}
	return nil
}

func (r *UserRepository) Find(nickname string) (*models.User, error) {
	user := new(models.User)
	 err := r.db.QueryRow(
		"SELECT about, email, fullname, nickname " +
			"FROM users " +
			"WHERE nickname = $1 ",
			nickname,
	).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusNotFound
		return nil, rerr
	}
	return user, nil
}

func (r *UserRepository) IsExist(nickname string, email string) bool {
	row := r.db.QueryRow(
		"SELECT email, nickname " +
			"FROM users " +
			"WHERE nickname = $1 AND email = $2",
		nickname,
		email,
	)
	if row == nil {
		return false
	}

	return true
}

func (r *UserRepository) Update(user *models.User) error {

}

func (r *UserRepository) CheckEmail (email string) bool {

}