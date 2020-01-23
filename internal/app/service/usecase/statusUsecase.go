package usecase

import (
	"database/sql"
	forumU "github.com/Andronovdima/tpark-db-forum/internal/app/forum/usecase"
	post "github.com/Andronovdima/tpark-db-forum/internal/app/post/usecase"
	threadU "github.com/Andronovdima/tpark-db-forum/internal/app/thread/usecase"
	userU "github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"log"
)

type ServiceUsecase struct {
	threadUC *threadU.ThreadUsecase
	forumUC  *forumU.ForumUsecase
	userUC   *userU.UserUsecase
	postUC   *post.PostUsecase
	db       *sql.DB
}

func NewServiceUsecase(tr *threadU.ThreadUsecase, fr *forumU.ForumUsecase, ur *userU.UserUsecase, pr *post.PostUsecase, thisdb *sql.DB) *ServiceUsecase {
	serviceUsecase := &ServiceUsecase{
		threadUC: tr,
		forumUC:  fr,
		userUC:   ur,
		postUC:   pr,
		db:       thisdb,
	}
	return serviceUsecase
}

func (s *ServiceUsecase) GetStatus() *models.Status {
	st := new(models.Status)
	st.Forum = s.forumUC.Status()
	st.Thread = s.threadUC.Status()
	st.User = s.userUC.Status()
	st.Post = s.postUC.Status()
	return st
}

func (s *ServiceUsecase) ClearData() {
	_ = truncateData(s.db)
}

func truncateData(db *sql.DB) error {
	query := `TRUNCATE TABLE forums CASCADE;
	TRUNCATE TABLE threads CASCADE;
	TRUNCATE TABLE users CASCADE;
	TRUNCATE TABLE posts CASCADE;
	TRUNCATE TABLE votes CASCADE;`
	if _, err := db.Exec(query); err != nil {
		log.Println(err.Error())
	}
	return nil
}
