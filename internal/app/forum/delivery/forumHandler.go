package delivery

import (
	"encoding/json"
	forum "github.com/Andronovdima/tpark-db-forum/internal/app/forum/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/app/respond"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type ForumHandler struct {
	ForumUsecase forum.ForumUsecase
	logger         *zap.SugaredLogger
	sessionStore   sessions.Store
}

func NewForumHandler(m *mux.Router, uc forum.ForumUsecase, logger *zap.SugaredLogger) {
	handler := &ForumHandler{
		ForumUsecase: 	uc,
		logger:         logger,
	}
	//
	m.HandleFunc("/", handler.HandleHello)
	m.HandleFunc("/forum/create", handler.HandleCreateForum).Methods(http.MethodPost)
	m.HandleFunc("/forum/{slug}/details", handler.HandleGetForum).Methods(http.MethodGet)
}

func (f *ForumHandler) HandleHello (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respond.Respond(w, r, http.StatusOK, "hello from server")
}

func (f *ForumHandler) HandleCreateForum (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()

	thisForum := new(models.Forum)

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(thisForum)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateForum<-Decode: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	fr , err := f.ForumUsecase.CreateForum(thisForum)
	if err != nil {
		rerr := err.(*models.HttpError)
		if rerr.StatusCode == http.StatusConflict {
			respond.Respond(w, r, http.StatusConflict, fr)
			return
		}
		err = errors.Wrapf(err, "HandleCreateForum<-CreateForum: ")
		respond.Error(w, r, rerr.StatusCode , rerr)
		return
	}

	respond.Respond(w, r, http.StatusCreated, fr)
	return
}

func (f *ForumHandler) HandleGetForum (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)
	slug := vars["slug"]
	fr , err := f.ForumUsecase.Find(slug)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, fr)
	return
}