package delivery

import (
	"github.com/Andronovdima/tpark-db-forum/internal/app/respond"
	service "github.com/Andronovdima/tpark-db-forum/internal/app/service/usecase"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type ServiceHandler struct {
	serviceUC  service.ServiceUsecase
	logger       *zap.SugaredLogger
}

func NewServiceHandler(m *mux.Router, uc service.ServiceUsecase, logger *zap.SugaredLogger) {
	handler := &ServiceHandler{
		serviceUC: uc,
		logger:       logger,
	}

	m.HandleFunc("/service/status", handler.HandleStatus).Methods(http.MethodGet)
	m.HandleFunc("/service/clear", handler.HandleClearData).Methods(http.MethodPost)
}

func (s *ServiceHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleStatus<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()
	status := s.serviceUC.GetStatus()

	respond.Respond(w, r, http.StatusOK, status)
}

func (s *ServiceHandler) HandleClearData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleClearData<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()

	s.serviceUC.ClearData()
	respond.Respond(w, r, http.StatusOK, nil)
}
