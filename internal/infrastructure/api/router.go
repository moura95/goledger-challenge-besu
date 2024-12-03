package api

import (
	"github.com/gin-gonic/gin"
	"github.com/moura95/goledger-challenge-besu/internal/application/service"
	"go.uber.org/zap"
)

type ContractRouter interface {
	SetupContractRoute(routers *gin.RouterGroup)
}

type StorageAPI struct {
	service service.StorageService
	logger  *zap.SugaredLogger
}

func NewReceiverRouter(s *service.StorageService, log *zap.SugaredLogger) *StorageAPI {
	return &StorageAPI{
		service: *s,
		logger:  log,
	}
}

func (r *StorageAPI) SetupContractRoute(routers *gin.RouterGroup) {
	routers.GET("/contract", r.get)
	routers.POST("/Contract", r.set)

}
