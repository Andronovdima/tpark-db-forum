package delivery

import (
	"encoding/json"
	"github.com/Andronovdima/tpark-db-forum/internal/app/respond"
	"github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)
type UserHandler struct {
	UserUsecase usecase.UserUsecase
	logger         *zap.SugaredLogger
	sessionStore   sessions.Store
}

func NewUserHandler(m *mux.Router, uc usecase.UserUsecase, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &UserHandler{
		UserUsecase: 	uc,
		logger:         logger,
		sessionStore:   sessionStore,
	}

	m.HandleFunc("/user/{nickname}/create", handler.HandleCreateUser).Methods(http.MethodPost)
	m.HandleFunc("/user/{nickname}/profile", handler.HandleGetProfile).Methods(http.MethodGet)
	m.HandleFunc("/user/{nickname}/profile", handler.HandleGetProfile).Methods(http.MethodPost)
}
func (u *UserHandler) HandleCreateUser (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	thisUser := new(models.User)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(thisUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser:")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	us , err := u.UserUsecase.CreateUser(thisUser, nickname)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	respond.Respond(w, r, http.StatusCreated, us)
}

func (u *UserHandler) HandleGetProfile (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	us, err := u.UserUsecase.Find(nickname)
	if err != nil {
		nerr := err.(*models.HttpError)
		respond.Error(w, r, nerr.StatusCode, nerr)
	}

	respond.Respond(w, r, http.StatusOK, us)
}

func (u *UserHandler) HandleUpdateProfile (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	thisUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(thisUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleUpdateProfile:")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	us , err := u.UserUsecase.UpdateProfile(nickname, thisUser)
	if err != nil {
		nerr := err.(*models.HttpError)
		respond.Error(w, r, nerr.StatusCode, nerr)
	}

	respond.Respond(w, r, http.StatusOK, us)
}