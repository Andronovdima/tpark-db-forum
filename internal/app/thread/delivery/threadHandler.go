package delivery

import (
	"encoding/json"
	"github.com/Andronovdima/tpark-db-forum/internal/app/respond"
	thread "github.com/Andronovdima/tpark-db-forum/internal/app/thread/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ThreadHandler struct {
	ThreadUsecase thread.ThreadUsecase
	logger        *zap.SugaredLogger
}

func NewForumHandler(m *mux.Router, uc thread.ThreadUsecase, logger *zap.SugaredLogger) {
	handler := &ThreadHandler{
		ThreadUsecase: uc,
		logger:        logger,
	}

	m.HandleFunc("/forum/{slug}/create", handler.HandleCreateThread).Methods(http.MethodPost)
	m.HandleFunc("/forum/{slug}/threads", handler.HandleGetThreads).Methods(http.MethodGet)
	m.HandleFunc("/thread/{slug_or_id}/details", handler.HandleGetThread).Methods(http.MethodGet)
	m.HandleFunc("/thread/{slug_or_id}/details", handler.HandleEditThread).Methods(http.MethodPost)
	m.HandleFunc("/thread/{slug_or_id}/vote", handler.HandleVoteThread).Methods(http.MethodPost)
}

func (t *ThreadHandler) HandleCreateThread(w http.ResponseWriter, r *http.Request) {
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
		err = errors.Wrapf(err, "HandleCreateSlug")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	th, err := t.ThreadUsecase.CreateThread(thisThread, slug)
	if err != nil {
		rerr := err.(*models.HttpError)
		if rerr.StatusCode == http.StatusConflict{
			respond.Respond(w, r, http.StatusConflict, th)
			return
		}

		respond.Error(w, r, rerr.StatusCode , rerr)
		return
	}

	respond.Respond(w, r, http.StatusCreated, th)
	return
}

func (t *ThreadHandler) HandleGetThreads(w http.ResponseWriter, r *http.Request) {
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

	threads, err := t.ThreadUsecase.GetThreads(slug, limit, since, desc)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, threads)
	return
}

func (t *ThreadHandler) HandleGetThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleGetThread<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)
	slugq := vars["slug_or_id"]

	var thr *models.Thread

	slugID, err := strconv.Atoi(slugq)
	if err != nil {
		isExist := t.ThreadUsecase.IsExistBySlug(slugq)
		if !isExist {
			respond.Error(w, r, http.StatusNotFound, errors.New("thread with this id doesn't exist"))
			return
		}

		ID := t.ThreadUsecase.GetIDBySlug(slugq)

		thr, err = t.ThreadUsecase.GetThreadByID(int32(ID))
		if err != nil {
			rerr := err.(*models.HttpError)
			respond.Error(w, r, rerr.StatusCode, rerr)
			return
		}

	} else {
		isExist := t.ThreadUsecase.IsExistByID(slugID)
		if !isExist {
			respond.Error(w, r, http.StatusNotFound, errors.New("thread with this id doesn't exist"))
			return
		}

		thr, err = t.ThreadUsecase.GetThreadByID(int32(slugID))
		if err != nil {
			rerr := err.(*models.HttpError)
			respond.Error(w, r, rerr.StatusCode, rerr)
			return
		}
	}

	respond.Respond(w, r, http.StatusOK, thr)
}

func (t *ThreadHandler) HandleEditThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditThread<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)

	slugq := vars["slug_or_id"]

	thrUpd := new(models.ThreadUpdate)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(thrUpd)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditThread")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	thr, err := t.ThreadUsecase.EditThread(slugq, thrUpd)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, thr)
	return
}

func (t *ThreadHandler) HandleVoteThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]

	vote := new(models.Vote)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(vote)
	if err != nil {
		err = errors.Wrapf(err, "HandleVoteThread")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	th, err := t.ThreadUsecase.VoteThread(slug, vote)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode , rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, th)
	return
}