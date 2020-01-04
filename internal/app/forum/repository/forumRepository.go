package repository

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
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
