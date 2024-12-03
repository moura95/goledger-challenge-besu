package server

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/moura95/goledger-challenge-besu/config"
	"github.com/moura95/goledger-challenge-besu/docs"
	"github.com/moura95/goledger-challenge-besu/internal/infrastructure/api"
	"github.com/moura95/goledger-challenge-besu/internal/infrastructure/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

type Server struct {
	store  *sqlx.DB
	router *gin.Engine
	config *config.Config
	logger *zap.SugaredLogger
}

//      @title                  Besu Api
//      @version                1.0
//      @description    Api for Interact With Smart Contract
//      @termsOfService http://swagger.io/terms/

//      @license.name   Apache 2.0
//      @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host           localhost:8080
// @BasePath       /api/v1
func NewServer(cfg config.Config, store *sqlx.DB, log *zap.SugaredLogger) *Server {

	server := &Server{
		store:  store,
		config: &cfg,
		logger: log,
	}
	var router *gin.Engine
	docs.SwaggerInfo.BasePath = ""

	router = gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Middleware Rate Limiter
	router.Use(middleware.RateLimitMiddleware())
	// Gzip
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Middleware Cors
	router.Use(middleware.CORSMiddleware())

	// Init all Routers
	api.CreateRoutesV1(store, server.config, router, log)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func RunGinServer(cfg config.Config, store *sqlx.DB, log *zap.SugaredLogger) {
	server := NewServer(cfg, store, log)

	_ = server.Start(cfg.HTTPServerAddress)
}
