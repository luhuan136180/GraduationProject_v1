package apiserver

import (
	"context"
	"github.com/robfig/cron/v3"
	"net/http"
	"v1/pkg/apis/v1/auth"
	"v1/pkg/apis/v1/contract"
	"v1/pkg/apis/v1/interview"
	"v1/pkg/apis/v1/project"
	"v1/pkg/apis/v1/resume"
	"v1/pkg/apis/v1/system"
	"v1/pkg/apiserver/imsystem"
	contractServe "v1/pkg/contract"

	"v1/pkg/apiserver/middleware"
	"v1/pkg/client/cache"
	"v1/pkg/logger"
	"v1/pkg/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type APIServer struct {
	Server  *http.Server
	Crontab *cron.Cron

	// webservice Router, where all RESTful API defines
	router *gin.Engine

	TokenManager token.Manager

	// mysql client
	RDBClient *gorm.DB

	// chat server
	ChatServer *imsystem.Server

	CacheClient cache.Interface

	//
	BlockClint *contractServe.Client
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.router = gin.New()
	s.router.ContextWithFallback = true
	s.router.Use(gin.Recovery())
	s.router.Use(logger.GinLogger())
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	}))
	s.router.Use(middleware.WithLanguage())

	if err := s.initSystem(); err != nil {
		zap.L().Panic("init system failed", zap.Error(err))
	}

	s.installAPIs()
	s.Server.Handler = s.router
	return nil
}

func (s *APIServer) Run(stopCh <-chan struct{}) (err error) {
	go func() {
		<-stopCh
		_ = s.Server.Shutdown(context.Background())
		s.Crontab.Stop()
	}()

	s.Crontab.Start()

	zap.L().Info("Start listening", zap.String("addr", s.Server.Addr))
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}

	return err
}

// add API Group
func (s *APIServer) installAPIs() {
	apiV1Group := s.router.Group("/api/v1")
	apiV1Group.Use(middleware.AddAuditLog(s.RDBClient))
	auth.RegisterRouter(apiV1Group, s.TokenManager, s.CacheClient, s.RDBClient)
	system.RegisterRouter(apiV1Group, s.TokenManager, s.CacheClient, s.RDBClient)
	project.RegisterRouter(apiV1Group, s.TokenManager, s.CacheClient, s.RDBClient)
	resume.RegisterRouter(apiV1Group, s.TokenManager, s.CacheClient, s.RDBClient)
	interview.RegisterRouter(apiV1Group, s.TokenManager, s.CacheClient, s.RDBClient)
	contract.RegisterRouter(apiV1Group, s.TokenManager, s.CacheClient, s.RDBClient, s.BlockClint)
	// ai.RegisterRouter(apiV1Group, s.TokenManager, s.CacheClient, s.RDBClient)
}
