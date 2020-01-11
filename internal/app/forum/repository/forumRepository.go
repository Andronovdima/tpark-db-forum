package repository

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(thisDB *sql.DB) *ForumRepository {
	forumRep := &ForumRepository{
		db: thisDB,
	}
	return forumRep
}

func (forumRep *ForumRepository) Create(forum *models.Forum) error {
	_ , err := forumRep.db.Exec(
		"INSERT INTO forums slug, title, user, posts, threads " +
			"VALUES $1, $2, $3, $4, $5",
			forum.Slug,
			forum.Title,
			forum.User,
			0,
			0,
	)
	return err
}

func (forumRep *ForumRepository) IsExist(slug string) bool {
	row := forumRep.db.QueryRow(
		"SELECT slug " +
			"FROM forums " +
			"WHERE slug = $1",
		slug,
	)
	if row == nil {
		return false
	}

	return true

}

func (forumRep *ForumRepository) Find(slug string) (*models.Forum, error) {
	forum := new(models.Forum)
	err := forumRep.db.QueryRow(
		"SELECT posts, slug, threads, title, user " +
			"FROM forums " +
			"WHERE slug = $1 ",
		slug,
	).Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusNotFound
		return nil, rerr
	}
	return forum, nil
}