package internal

import (
	"context"
	"github.com/pengcainiao/pengcainiao/okr/internal/middleware"
	"github.com/pengcainiao/pengcainiao/okr/internal/v1/api"
	"github.com/pengcainiao/zero/core/logx"
	"github.com/pengcainiao/zero/rest"
	"github.com/pengcainiao/zero/rest/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Setup 开启服务
// @title 飞项核心业务API文档
// @version 1.0
// @description 飞项核心业务API文档
// @schemes http https
// @host 127.0.0.1:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Setup() {
	//setupGRPCServer()
	//setupTimerExec()
	gracefullShutdown(setupHTTPServer())
}

func setupHTTPServer() *http.Server {

	router := rest.NewGinServer()
	router.Use(httprouter.Recovery())
	router.Use(middleware.Cors())

	v1 := router.Group("/v1")
	{
		var (
			objective api.ObjectiveController
		)

		v1.GET("test", objective.First)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,
		IdleTimeout:  time.Second * 30,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen 8080 失败")
		}
	}()
	return srv
}

/**
 * 启动GRPC服务
 * @synopsis setupGRPCServer
 * @return
 */
//func setupGRPCServer() {
//	go func() {
//		if err := grpcbase.RegisterServer(grpcsvc.NewService()); err != nil {
//			log.Fatal().Err(err).Msg("listen 8084 失败")
//		}
//	}()
//}
/**
 * 启动健康检查
 * @synopsis setupHealthZ
 * @return
 */
func gracefullShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // nolint
	<-quit
	logx.NewTraceLogger(context.Background()).Debug().Msg("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logx.NewTraceLogger(context.Background()).Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logx.NewTraceLogger(context.Background()).Debug().Msg("Server exiting")
}
