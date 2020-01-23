package delivery

import (
	"encoding/json"
	post "github.com/Andronovdima/tpark-db-forum/internal/app/post/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/app/respond"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	PostUsecase  post.PostUsecase
	logger       *zap.SugaredLogger
}

func NewPostHandler(m *mux.Router, uc post.PostUsecase, logger *zap.SugaredLogger) {
	handler := &PostHandler{
		PostUsecase: uc,
		logger:       logger,
	}

	m.HandleFunc("/thread/{slug_or_id}/create", handler.HandleCreatePosts).Methods(http.MethodPost)
	m.HandleFunc("/post/{id}/details", handler.HandleEditPost).Methods(http.MethodPost)
	m.HandleFunc("/post/{id}/details", handler.HandleGetPost).Methods(http.MethodGet)
	m.HandleFunc("/thread/{slug_or_id}/posts", handler.HandleGetPosts).Methods(http.MethodGet)
	//m.HandleFunc("/thread/{slug_or_id}/posts", handler.HandleGetThreadPosts).Methods(http.MethodGet)
}

func (p *PostHandler) HandleCreatePosts (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreatePosts<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)
	slug := vars["slug_or_id"]

	posts := new([]models.Post)

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(posts)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreatePosts<-Decode: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	posts, err = p.PostUsecase.CreatePosts(slug, posts)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusCreated, posts)
	return
}

func (p *PostHandler) HandleEditPost (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditPost<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err:= strconv.ParseInt(ids, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPost:")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	postUpdate := new(models.PostUpdate)

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(postUpdate)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPost:")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	thisPost, err := p.PostUsecase.EditPost(id, postUpdate)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, thisPost)
}


func (p *PostHandler) HandleGetPost (w http.ResponseWriter, r *http.Request) {
	var includeForum, includeUser, includeThread bool
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleGetPost<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)
	ids := vars["id"]
	id , err:= strconv.ParseInt(ids, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPost:")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	relatedq := r.URL.Query()["related"]
	if relatedq != nil {
		related := strings.Split(relatedq[0], ",")

		if Contains(related, "user") {
			includeUser = true
		}

		if Contains(related, "forum") {
			includeForum = true
		}

		if Contains(related, "thread") {
			includeThread = true
		}
	}

	postFull, err := p.PostUsecase.GetInfo(id, includeForum, includeThread, includeUser)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, postFull)
}

func Contains(str []string, s string) bool {
	for _, i := range str {
		if i == s {
			return true
		}
	}
	return false
}

func (p *PostHandler) HandleGetPosts (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleGetPosts<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	vars := mux.Vars(r)
	slugID := vars["slug_or_id"]

	limit := 100
	sort := "flat"
	since := -1
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
		since, err = strconv.Atoi(sinceq[0])
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.New("Unavailable type of limit"))
			return
		}
	}

	descq, ok := r.URL.Query()["desc"]
	if ok && len(descq[0]) > 0 {
		desc, err = strconv.ParseBool(descq[0])
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.New("Unavailable type of desc"))
			return
		}
	}

	sortq, ok := r.URL.Query()["sort"]
	if ok && len(sortq[0]) > 0 {
		if sortq[0] != "flat" && sortq[0] != "tree" && sortq[0] != "parent_tree"{
			respond.Error(w, r, http.StatusBadRequest, errors.New("Unavailable value of sort"))
			return
		}
		sort = sortq[0]
	}


	posts, err := p.PostUsecase.GetPosts(slugID, limit, since, sort, desc)
	if err != nil {
		rerr := err.(*models.HttpError)
		respond.Error(w, r, rerr.StatusCode, rerr)
		return
	}

	respond.Respond(w, r, http.StatusOK, posts)
	return
}


//func (p *PostHandler) HandleGetThreadPosts(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	vars := mux.Vars(r)
//	slugOrId := vars["slug_or_id"]
//
//	posts, err := p.PostUsecase.GetThreadPosts(slugOrId, r.URL.Query())
//
//	if err != nil {
//		respond.Error(w, r, http.StatusNotFound, err)
//		return
//	}
//
//	respond.Respond(w, r, http.StatusOK, &posts)
//}