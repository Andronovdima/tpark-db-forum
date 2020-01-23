package repository

import (
	"database/sql"
	"errors"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(thisDB *sql.DB) *PostRepository {
	postRep := &PostRepository{
		db: thisDB,
	}
	return postRep
}

func (pr *PostRepository) CreatePosts(threadID int32, forumSlug string, posts *[]models.Post) (int, error) {

	tx, err := pr.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	now := time.Now()

	queryStr := "INSERT INTO posts(id, parent, thread, forum, author, created, message, path) VALUES "
	values := []interface{}{}
	for _, post := range *posts {

		if post.Parent == 0 {
			queryStr += "(nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint]),"
			values = append(values, post.Parent, threadID, forumSlug, post.Author, now, post.Message)
		} else {
			var parentThreadId int32
			err = pr.db.QueryRow("SELECT thread FROM posts WHERE id = $1",
				post.Parent,
			).Scan(&parentThreadId)
			if err != nil {
				_ = tx.Rollback()
				return 404, err
			}
			if parentThreadId != threadID {
				_ = tx.Rollback()
				return 409, errors.New("different Threads id in posts")
			}

			queryStr += " (nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"(SELECT path FROM posts WHERE id = ? AND thread = ?) || " +
				"currval(pg_get_serial_sequence('posts', 'id'))::bigint),"

			values = append(values, post.Parent, threadID, forumSlug, post.Author, now, post.Message, post.Parent, threadID)
		}


	}
	queryStr = strings.TrimSuffix(queryStr, ",")

	queryStr += " RETURNING  id, thread, forum, created "

	queryStr = ReplaceSQL(queryStr, "?")
	if len(*posts) > 0 {
		stmtButch, err := tx.Prepare(queryStr)
		if err != nil {
			return 404, err
		}
		rows, err := stmtButch.Query(values...)
		if err != nil {
			_ = tx.Rollback()
			return 404, err
		}
		i := 0
		for rows.Next() {
			err := rows.Scan(
				&(*posts)[i].Id,
				&(*posts)[i].Thread,
				&(*posts)[i].Forum,
				&(*posts)[i].Created,
			)
			i += 1

			if err != nil {
				tx.Rollback()
				return 500, err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return 500, err
	}

	//(*posts)[i].Forum = forumSlug
	//(*posts)[i].Thread = threadID
	//if (*posts)[i].Created != "" {
	//	err := pr.db.QueryRow(
	//		"INSERT INTO posts (author, created, forum, IsEdited, message, parent, thread) "+
	//			"VALUES ($1, $2, $3, $4, $5, $6, $7) " +
	//			"RETURNING id",
	//		(*posts)[i].Author,
	//		(*posts)[i].Created,
	//		(*posts)[i].Forum,
	//		false,
	//		(*posts)[i].Message,
	//		(*posts)[i].Parent,
	//		(*posts)[i].Thread,
	//	).Scan(&(*posts)[i].Id)
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	err := pr.db.QueryRow(
	//		"INSERT INTO posts (author,  forum, IsEdited, message, parent, thread) "+
	//			"VALUES ($1, $2, $3, $4, $5, $6) " +
	//			"RETURNING id",
	//		(*posts)[i].Author,
	//		(*posts)[i].Forum,
	//		false,
	//		(*posts)[i].Message,
	//		(*posts)[i].Parent,
	//		(*posts)[i].Thread,
	//	).Scan(&(*posts)[i].Id)
	//
	//	if err != nil {
	//		return err
	//	}
	//}
	//}
	return 200, nil
}

func (pr *PostRepository) IsExist(ID int64) bool {
	row := pr.db.QueryRow(
		"SELECT id "+
			"FROM posts "+
			"WHERE id = $1",
		ID,
	)
	if row.Scan(&ID) != nil {
		return false
	}

	return true
}

func (pr *PostRepository) EditPost(ID int64, update *models.PostUpdate) error {
	_, err := pr.db.Exec(
		"UPDATE posts "+
			"SET message = $1 , isEdited = $2 "+
			"WHERE id = $3",
		update.Message,
		true,
		ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PostRepository) Find(ID int64) (*models.Post, error) {
	var p models.Post
	err := pr.db.QueryRow(
		"SELECT author, created, forum, IsEdited, message, parent, thread "+
			"FROM posts "+
			"WHERE id = $1",
		ID,
	).Scan(&p.Author, &p.Created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
	if err != nil {
		return nil, err
	}
	p.Id = ID
	return &p, nil
}

func (pr *PostRepository) Status() int64 {
	var count int64
	_ = pr.db.QueryRow(
		"SELECT COUNT(*) " +
			"FROM posts ",
	).Scan(&count)
	return count
}
//
func (pr *PostRepository) GetThreadPosts(ID int32, limit int , since int, sort string, desc bool) ([]*models.Post, error) {
	var rows *sql.Rows
	var err error
	posts := []*models.Post{}

	var descq, sign string
	sign = ">"
	descq = ""

	if desc {
		sign = "<"
		descq = "DESC"
	}

	switch sort {

	case "flat" :
		if since == -1 {
			q := "SELECT id, parent, thread, forum, author, created, message, isedited " +
				"FROM posts " +
				"WHERE thread = $1" +
				"ORDER BY created " + descq + " , id " + descq + " " +
				"LIMIT $2"

			rows , err = pr.db.Query(q,
				ID,
				limit,
			)
			if err != nil {
				return nil, err
			}
		} else {
			q := "SELECT id, parent, thread, forum, author, created, message, isedited " +
				"FROM posts " +
				"WHERE thread = $1 AND id " + sign + " $2 " +
				"ORDER BY created " + descq + " , id " + descq + " " +
				"LIMIT $3"

			rows , err = pr.db.Query(q,
				ID,
				since,
				limit,
			)
			if err != nil {
				return nil, err
			}
		}

	case "tree":
		if since == - 1 {
			q := "SELECT id, parent, thread, forum, author, created, message, isedited " +
				"FROM posts " +
				"WHERE thread = $1" +
				"ORDER BY path[1] " + descq + ", path " + descq + " " +
				"LIMIT $2"

			rows , err = pr.db.Query(q,
				ID,
				limit,
			)
			if err != nil {
				return nil, err
			}
		} else {
			q := "SELECT id, parent, thread, forum, author, created, message, isedited " +
				"FROM posts " +
				"WHERE thread = $1 AND path " + sign + " (SELECT path FROM posts WHERE id = $2) " +
				"ORDER BY path[1] " + descq + ", path " + descq + " " +
				"LIMIT $3"

			rows , err = pr.db.Query(q,
				ID,
				since,
				limit,
			)
			if err != nil {
				return nil, err
			}
		}

	case "parent_tree":

		if since == -1 {
			q := "SELECT id, parent, thread, forum, author, created, message, isedited " +
				"FROM posts " +
				"WHERE thread = $1 AND path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND parent = 0 " +
				"ORDER BY path[1] " + descq + ", path LIMIT $2)) " +
				"ORDER BY path[1] " + descq + ", path "

			rows , err = pr.db.Query(q,
				ID,
				limit,
			)
			if err != nil {
				return nil, err
			}
		} else {
			q := "SELECT id, parent, thread, forum, author, created, message, isedited " +
				"FROM posts " +
				"WHERE thread = $1 AND path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND parent = 0 " +
				"AND path " + sign + " (SELECT path[1:1] FROM posts WHERE id = $2) " +
				"ORDER BY path[1] " + descq + ", path LIMIT $3)) " +
				"ORDER BY path[1] " + descq + ", path "

			rows , err = pr.db.Query(q,
				ID,
				since,
				limit,
			)
			if err != nil {
				return nil, err
			}
		}

		}

	for rows.Next() {
		p := models.Post{}
		err := rows.Scan(&p.Id, &p.Parent, &p.Thread, &p.Forum, &p.Author, &p.Created, &p.Message, &p.IsEdited)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return posts, nil
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
