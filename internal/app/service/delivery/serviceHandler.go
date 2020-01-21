package delivery

import (
	service "github.com/Andronovdima/tpark-db-forum/internal/app/status/usecase"
	"github.com/gorilla/mux"
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

	m.HandleFunc("/service/{slug_or_id}/create", handler.HandleCreatePosts).Methods(http.MethodPost)
}