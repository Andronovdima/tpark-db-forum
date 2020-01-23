package apiserver

import (
	"database/sql"
	forumH "github.com/Andronovdima/tpark-db-forum/internal/app/forum/delivery"
	forumR "github.com/Andronovdima/tpark-db-forum/internal/app/forum/repository"
	forumU "github.com/Andronovdima/tpark-db-forum/internal/app/forum/usecase"
	postH "github.com/Andronovdima/tpark-db-forum/internal/app/post/delivery"
	postR "github.com/Andronovdima/tpark-db-forum/internal/app/post/repository"
	postU "github.com/Andronovdima/tpark-db-forum/internal/app/post/usecase"
	serviceH "github.com/Andronovdima/tpark-db-forum/internal/app/service/delivery"
	serviceU "github.com/Andronovdima/tpark-db-forum/internal/app/service/usecase"
	threadH "github.com/Andronovdima/tpark-db-forum/internal/app/thread/delivery"
	threadR "github.com/Andronovdima/tpark-db-forum/internal/app/thread/repository"
	threadU "github.com/Andronovdima/tpark-db-forum/internal/app/thread/usecase"
	userH "github.com/Andronovdima/tpark-db-forum/internal/app/user/delivery"
	userR "github.com/Andronovdima/tpark-db-forum/internal/app/user/repository"
	userU "github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Mux          *mux.Router
	Config       *Config
	Logger       *zap.SugaredLogger
}

func NewServer(config *Config, logger *zap.SugaredLogger) (*Server, error) {
	s := &Server{
		Mux:          mux.NewRouter().PathPrefix("/api").Subrouter(),
		Logger:       logger,
		Config:       config,
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func (s *Server) ConfigureServer(db *sql.DB) {

	userRep := userR.NewUserRepository(db)
	forumRep := forumR.NewForumRepository(db)
	threadRep := threadR.NewThreadRepository(db)
	postRep := postR.NewPostRepository(db)


	userUc := userU.NewUserUsecase(userRep, forumRep)
	forumUc := forumU.NewForumUsecase(forumRep, userUc)
	threadUc := threadU.NewUserUsecase(threadRep, userUc, forumUc)
	postUc := postU.NewPostUsecase(postRep, threadUc, forumUc, userUc)
	serviceUc := serviceU.NewServiceUsecase(threadUc, forumUc, userUc, postUc, db)

	userH.NewUserHandler(s.Mux, *userUc, s.Logger)
	forumH.NewForumHandler(s.Mux, *forumUc, s.Logger)
	threadH.NewForumHandler(s.Mux, *threadUc, s.Logger)
	postH.NewPostHandler(s.Mux, *postUc, s.Logger)
	serviceH.NewServiceHandler(s.Mux, *serviceUc, s.Logger)

}
