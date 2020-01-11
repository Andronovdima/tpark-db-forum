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

func NewForumHandler(m *mux.Router, uc forum.ForumUsecase, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &ForumHandler{
		ForumUsecase: 	uc,
		logger:         logger,
		sessionStore:   sessionStore,
	}
	//
	m.HandleFunc("/forum/create", handler.HandleCreateForum).Methods(http.MethodPost)
	m.HandleFunc("/forum/{slug}/create", handler.HandleCreateSlug).Methods(http.MethodPost)
}

func (f *ForumHandler) HandleCreateForum (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
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

	createdForum , err := f.ForumUsecase.CreateForum(thisForum)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateForum<-Decode: ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	respond.Respond(w, r, http.StatusCreated, createdForum)
}

func (f *ForumHandler) HandleCreateSlug (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()
	vars := mux.Vars(r)
	slug := vars["slug"]
	thisThread := new(models.Thread)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(thisThread)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateSlug:")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	th , err := f.ForumUsecase.CreateThread(thisThread, slug)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	respond.Respond(w, r, http.StatusCreated, th)
}
