package delivery

import (
	"encoding/json"
	"github.com/Andronovdima/tpark-db-forum/internal/app/respond"
	thread "github.com/Andronovdima/tpark-db-forum/internal/app/thread/usecase"
	"github.com/Andronovdima/tpark-db-forum/internal/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type ThreadHandler struct {
	ThreadUsecase thread.ThreadUsecase
	logger        *zap.SugaredLogger
	sessionStore  sessions.Store
}

func NewForumHandler(m *mux.Router, uc thread.ThreadUsecase, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &ThreadHandler{
		ThreadUsecase: uc,
		logger:        logger,
		sessionStore:  sessionStore,
	}

	m.HandleFunc("/forum/{slug}/create", handler.HandleCreateSlug).Methods(http.MethodPost)
}

func (t *ThreadHandler) HandleCreateSlug(w http.ResponseWriter, r *http.Request) {
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
	th, err := t.ThreadUsecase.CreateThread(thisThread, slug)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	respond.Respond(w, r, http.StatusCreated, th)
}
