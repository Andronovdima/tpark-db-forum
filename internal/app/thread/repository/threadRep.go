package repository

import (
	"database/sql"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
	"strings"
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


func (tr *ThreadRepository) Find(slug string) (*models.Thread, error) {
	thread := new(models.Thread)
	err := tr.db.QueryRow(
		"SELECT author, created, id, forum, message, slug, title, votes " +
			"FROM threads " +
			"WHERE LOWER(slug) = $1 ",
		strings.ToLower(slug),
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
	thread.Forum = forumSlug
	if thread.Created == "" {
		err := tr.db.QueryRow(
			"INSERT INTO threads (author, forum, message, slug, title, votes)"+
				"VALUES ($1, $2, $3, $4, $5, $6) "+
				"RETURNING id",
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Slug,
			thread.Title,
			0,
		).Scan(&thread.Id)
		return err
	} else {
		err := tr.db.QueryRow(
			"INSERT INTO threads (author, created, forum, message, slug, title, votes)"+
				"VALUES ($1, $2, $3, $4, $5, $6, $7) "+
				"RETURNING id",
			thread.Author,
			thread.Created,
			thread.Forum,
			thread.Message,
			thread.Slug,
			thread.Title,
			0,
		).Scan(&thread.Id)
		return err
	}
}

func (tr *ThreadRepository) GetThreads(forumSlug string, limit int, created string, desc bool) ([]models.Thread, error) {
	var rows *sql.Rows
	var err error
	threads := []models.Thread{}

	var descq, sign string
	sign = ">="
	descq = ""

	if desc {
		sign = "<="
		descq = "DESC"
	}

	if created == "" {
		q := "SELECT author, created, id, forum, message, slug, title, votes " +
			"FROM threads " +
			"WHERE LOWER(forum) = $1" +
			"ORDER BY created " + descq + " " +
			"LIMIT $2"

		rows , err = tr.db.Query(q,
			strings.ToLower(forumSlug),
			limit,
		)
		if err != nil {
			return nil, err
		}
	} else {
		q := "SELECT author, created, id, forum, message, slug, title, votes " +
			"FROM threads " +
			"WHERE LOWER(forum) = $1 AND created " + sign + " $2 " +
			"ORDER BY created " + descq + " " +
			"LIMIT $3"

		rows , err = tr.db.Query(q,
			strings.ToLower(forumSlug),
			created,
			limit,
		)
		if err != nil {
			return nil, err
		}
	}


	for rows.Next() {
		t := models.Thread{}
		err := rows.Scan(&t.Author, &t.Created, &t.Id, &t.Forum, &t.Message, &t.Slug, &t.Title, &t.Votes)
		if err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}


func (tr *ThreadRepository) IsExist(th *models.Thread) (bool, int32) {
	var id int32
	row := tr.db.QueryRow(
		"SELECT id " +
			"FROM threads " +
			"WHERE (slug <> '' AND LOWER(slug) = $1) OR (LOWER(author) = $2 AND title = $3 AND message = $4 AND LOWER(forum) = $5)",
		strings.ToLower(th.Slug),
		strings.ToLower(th.Author),
		th.Title,
		th.Message,
		strings.ToLower(th.Forum),
	)
	if row.Scan(&id) != nil {
		return false, 0
	}

	return true, id
}


func (tr *ThreadRepository) IsExistByID(ID int) bool {
	row := tr.db.QueryRow(
		"SELECT id " +
			"FROM threads " +
			"WHERE id = $1",
		ID,
	)
	if row.Scan(&ID) != nil {
		return false
	}

	return true
}


func (tr *ThreadRepository) GetIDBySlug(slug string) int {
	var ID int
	row := tr.db.QueryRow(
		"SELECT id " +
			"FROM threads " +
			"WHERE LOWER(slug) = $1",
		strings.ToLower(slug),
	)
	if row.Scan(&ID) != nil {
		return 0
	}

	return ID
}

func (tr *ThreadRepository) GetForum(slug string, id int32) (string, error) {
	var forum string
	row := tr.db.QueryRow(
		"SELECT forum " +
			"FROM threads " +
			"WHERE (slug <> '' AND LOWER(slug) = $1) OR id = $2",
		strings.ToLower(slug),
		id,
	)
	err := row.Scan(&forum)
	if err != nil {
		return "", err
	}

	return forum, nil
}


func (tr *ThreadRepository) FindByID(id int32) (*models.Thread, error) {
	thread := new(models.Thread)
	err := tr.db.QueryRow(
		"SELECT author, created, id, forum, message, slug, title, votes " +
			"FROM threads " +
			"WHERE id = $1 ",
		id,
	).Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		rerr := new(models.HttpError)
		rerr.StringErr = err.Error()
		rerr.StatusCode = http.StatusInternalServerError
		return nil, rerr
	}
	return thread, nil
}

func (tr *ThreadRepository) EditThread(ID int32, thrUpd *models.ThreadUpdate) error {
	_, err := tr.db.Exec(
		"UPDATE threads " +
			"SET message = $1 , title = $2 " +
			"WHERE id = $3",
		thrUpd.Message,
		thrUpd.Title,
		ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) IsVoted(nickname string, threadID int32) bool {
	row := tr.db.QueryRow(
		"SELECT nickname " +
			"FROM votes " +
			"WHERE thread = $1 AND nickname = $2",
		threadID,
		nickname,
	)
	if row.Scan(&nickname) != nil {
		return false
	}

	return true
}


func (tr *ThreadRepository) IsVotedSame(vote *models.Vote, threadID int32) bool {
	var str string
	row := tr.db.QueryRow(
		"SELECT nickname " +
			"FROM votes " +
			"WHERE thread = $1 AND nickname = $2 AND voice = $3",
		threadID,
		vote.Nickname,
		vote.Voice,
	)
	if row.Scan(&str) != nil {
		return false
	}

	return true
}


func (tr *ThreadRepository) CreateVote(threadID int32, vote *models.Vote) error {
	_, err := tr.db.Exec(
		"INSERT INTO votes (thread, nickname, voice)" +
			"VALUES ($1, $2, $3) ",
			threadID,
			vote.Nickname,
			vote.Voice,
	)
	return err
}

func (tr *ThreadRepository) AddVote(threadID int32, count int) error {
	_, err := tr.db.Exec(
		"UPDATE threads " +
			"SET votes = votes + $1 " +
			"WHERE id = $2",
			count,
			threadID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) EditVote(threadID int32, vote *models.Vote) error {
	_, err := tr.db.Exec(
		"UPDATE votes " +
			"SET voice = $1" +
			"WHERE thread = $2 AND nickname = $3",
		vote.Voice,
		threadID,
		vote.Nickname,
	)
	if err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) Status() int32 {
	var count int32
	_ = tr.db.QueryRow(
		"SELECT COUNT(*) " +
			"FROM threads ",
	).Scan(&count)
	return count
}

func (tr *ThreadRepository) IsExistBySlug(slug string) bool {
	row := tr.db.QueryRow(
		"SELECT id "+
			"FROM threads "+
			"WHERE slug <> '' AND LOWER(slug) = $1",
		strings.ToLower(slug),
	)
	if row.Scan(&slug) != nil {
		return false
	}

	return true
}