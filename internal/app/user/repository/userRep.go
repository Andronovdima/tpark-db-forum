package repository

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
	"strings"
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
		"INSERT INTO users (about, email, fullname, nickname)" +
			"VALUES ($1, $2, $3, $4)",
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
			"WHERE LOWER(nickname) = $1 ",
			strings.ToLower(nickname),
	).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusNotFound
		return nil, rerr
	}
	return user, nil
}

func (r *UserRepository) IsExist(nickname string) bool {
	row := r.db.QueryRow(
		"SELECT nickname " +
			"FROM users " +
			"WHERE LOWER(nickname) = $1",
		strings.ToLower(nickname),
	)

	if row.Scan(&nickname) != nil {
		return false
	}

	return true
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.db.Exec(
		"UPDATE users " +
			"SET about = $1 , email = $2, fullname = $3 " +
			"WHERE nickname = $4",
			user.About,
			user.Email,
			user.Fullname,
			user.Nickname,
		)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = "Internal Error"
		rerr.StatusCode = http.StatusInternalServerError
		return rerr
	}

	return nil
}

//func (r *UserRepository) CheckEmail(email string) bool {
//	row := r.db.QueryRow(
//		"SELECT email " +
//			"FROM users " +
//			"WHERE email = $1",
//		email,
//	)
//	if row.Scan(&email) != nil {
//		return false
//	}
//
//	return true
//}

func (r *UserRepository) GetForumUsers(forumSlug string, limit int, since string, desc bool) ([]models.User, error) {
	users := []models.User{}
	var rows *sql.Rows
	var err error

	var descq, sign  string
	sign = ">"

	if desc {
		sign = "<"
		descq = "DESC"
	}
	if since == "" {

		q := "SELECT about, email, fullname, nickname FROM users " +
			"WHERE nickname IN (SELECT nickname FROM forum_users WHERE LOWER(forum) = LOWER($1)) " +
			"ORDER BY LOWER (nickname) " + descq + " LIMIT $2"
		rows, err = r.db.Query(
			q,
			forumSlug,
			limit,
			)
	} else {
		q := "SELECT about, email, fullname, nickname FROM users " +
			"WHERE nickname IN (SELECT nickname FROM forum_users WHERE LOWER(forum) = LOWER($1)) " +
			"AND LOWER(nickname) " + sign + " LOWER($2) " +
			"ORDER BY LOWER(nickname) " + descq + " LIMIT $3"
		rows, err = r.db.Query(
			q,
			forumSlug,
			since,
			limit,
		)
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Status() int32 {
	var count int32
	_ = r.db.QueryRow(
		"SELECT COUNT(*) " +
			"FROM users ",
	).Scan(&count)
	return count
}

func (r *UserRepository) IsExistEmail(email string) bool {
	row := r.db.QueryRow(
		"SELECT email " +
			"FROM users " +
			"WHERE Lower(email) = $1",
		strings.ToLower(email),
	)

	if row.Scan(&email) != nil {
		return false
	}

	return true
}


func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := new(models.User)
	err := r.db.QueryRow(
		"SELECT about, email, fullname, nickname " +
			"FROM users " +
			"WHERE LOWER(email) = $1 ",
		strings.ToLower(email),
	).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusNotFound
		return nil, rerr
	}

	return user, nil
}

func (r *UserRepository) GetRightNickname(nickname string) string {
	row := r.db.QueryRow(
		"SELECT nickname " +
			"FROM users " +
			"WHERE LOWER(nickname) = $1",
		strings.ToLower(nickname),
	)

	if row.Scan(&nickname) != nil {
		return ""
	}

	return nickname
}

