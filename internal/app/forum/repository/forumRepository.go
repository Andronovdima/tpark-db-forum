package repository

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
	"strings"
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
		"INSERT INTO forums (slug, title, author, posts, threads) " +
			"VALUES ($1, $2, $3, $4, $5)",
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
			"WHERE LOWER(slug) = $1",
		strings.ToLower(slug),
	)
	if row.Scan(&slug) != nil {
		return false
	}

	return true

}

func (forumRep *ForumRepository) Find(slug string) (*models.Forum, error) {
	forum := new(models.Forum)
	err := forumRep.db.QueryRow(
		"SELECT posts, slug, threads, title, author " +
			"FROM forums " +
			"WHERE LOWER(slug) = $1 ",
		strings.ToLower(slug),
	).Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusNotFound
		return nil, rerr
	}
	return forum, nil
}

func (forumRep *ForumRepository) AddThread(slug string) error {
	_ , err := forumRep.db.Exec(
		"UPDATE forums SET threads = threads + 1 " +
			"WHERE slug = $1",
		slug,
	)

	return err
}


func (forumRep *ForumRepository) AddPost(slug string, count int64) error {
	_ , err := forumRep.db.Exec(
		"UPDATE forums SET posts = posts + $1" +
			"WHERE slug = $2",
		count,
		slug,
	)

	return err
}


func (forumRep *ForumRepository) Status() int64 {
	var count int64
	_ = forumRep.db.QueryRow(
		"SELECT COUNT(*) " +
			"FROM forums ",
	).Scan(&count)
	return count
}

func (forumRep *ForumRepository) GetRightSlug(slug string) string {
	row := forumRep.db.QueryRow(
		"SELECT slug " +
			"FROM forums " +
			"WHERE LOWER(slug) = $1",
		strings.ToLower(slug),
	)

	if row.Scan(&slug) != nil {
		return ""
	}

	return slug
}

