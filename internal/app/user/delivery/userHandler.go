package delivery

import (
	"encoding/json"
	"github.com/Andronovdima/tpark-db-forum/internal/app/respond"
	"github.com/Andronovdima/tpark-db-forum/internal/app/user/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)
type UserHandler struct {
	UserUsecase usecase.UserUsecase
	logger         *zap.SugaredLogger
}

func NewUserHandler(m *mux.Router, uc usecase.UserUsecase, logger *zap.SugaredLogger) {
	handler := &UserHandler{
		UserUsecase: 	uc,
		logger:         logger,
	}

	m.HandleFunc("/user/{nickname}/create", handler.HandleCreateUser).Methods(http.MethodPost)
	m.HandleFunc("/user/{nickname}/profile", handler.HandleGetProfile).Methods(http.MethodGet)
	m.HandleFunc("/user/{nickname}/profile", handler.HandleUpdateProfile).Methods(http.MethodPost)
	m.HandleFunc("/forum/{slug}/users", 	handler.HandleGetForumUsers).Methods(http.MethodGet)
}
func (u *UserHandler) HandleCreateUser (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
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
		rerr:= err.(*models.HttpError)
		if rerr.StatusCode == http.StatusConflict {
			respond.Respond(w, r, http.StatusConflict, us)
			return
		}
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusCreated, (*us)[0])
	return
}

func (u *UserHandler) HandleGetProfile (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	us, err := u.UserUsecase.Find(nickname)
	if err != nil {
		nerr := err.(*models.HttpError)
		respond.Error(w, r, nerr.StatusCode, nerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, us)
}

func (u *UserHandler) HandleUpdateProfile (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
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
		return
	}

	respond.Respond(w, r, http.StatusOK, us)
}

func (u *UserHandler) HandleGetForumUsers (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()
	vars := mux.Vars(r)
	slug := vars["slug"]

	limit := 1000
	var since string
	var desc bool
	var err error

	limitq, ok := r.URL.Query()["limit"]
	if ok && len(limitq[0]) > 0 {
		limit, err = strconv.Atoi(limitq[0])
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.New("Unavailable type of limit"))
			return
		}
	}

	sinceq, ok := r.URL.Query()["since"]
	if ok && len(sinceq[0]) > 0 {
		since = sinceq[0]
	}


	descq, ok := r.URL.Query()["desc"]
	if ok && len(descq[0]) > 0 {
		desc, err = strconv.ParseBool(descq[0])
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.New("Unavailable type of desc"))
			return
		}
	}

	users, err := u.UserUsecase.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, users)
	return

}