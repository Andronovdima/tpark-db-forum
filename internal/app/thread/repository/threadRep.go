package repository

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
	"time"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(thisDB *sql.DB) *ThreadRepository {
	threadRep := &ThreadRepository{
		db: thisDB,
	}
	return threadRep
}

func (tr *ThreadRepository) IsExist(slug string) bool {
	row := tr.db.QueryRow(
		"SELECT slug " +
			"FROM threads " +
			"WHERE slug = $1",
		slug,
	)
	if row == nil {
		return false
	}

	return true
}


func (tr *ThreadRepository) Find(slug string) (*models.Thread, error) {
	thread := new(models.Thread)
	err := tr.db.QueryRow(
		"SELECT author, created, id, forum, message, slug, title, votes " +
			"FROM threads " +
			"WHERE slug = $1 ",
		slug,
	).Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}
	return thread, nil
}

func (tr *ThreadRepository) Create(thread *models.Thread, forumSlug string) error {
	err := tr.db.QueryRow(
		"INSERT INTO threads author, created, forum, message, slug, title, votes " +
			"VALUES $1, $2, $3, $4, $5, $6, $7 " +
			"RETURNING id",
		thread.Author,
		time.Now(),
		forumSlug,
		thread.Message,
		thread.Slug,
		thread.Title,
		0,
	).Scan(&thread.Id)
	return err
}
