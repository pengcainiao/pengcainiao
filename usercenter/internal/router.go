package internal

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/rest"
	"gitlab.com/a16624741591/zero/rest/httprouter"
	"gitlab.com/a16624741591/zero/rpcx/grpcbase"
	"net/http"
	"os"
	"os/signal"
	"pp/usercenter/internal/auth"
	"pp/usercenter/internal/auth/adapter/mysql"
	"pp/usercenter/internal/grpcsvc"
	"pp/usercenter/internal/middleware"
	v11 "pp/usercenter/internal/v1"
	"pp/usercenter/internal/v1/api"
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
	setupGRPCServer()
	//setupTimerExec()
	gracefullShutdown(setupHTTPServer())
}

func setupHTTPServer() *http.Server {
	env.RedisAddr = "127.0.0.1:6379"
	env.DbDSN = "penglonghui:Penglonghui!123!@tcp(119.29.5.54:3306)/okr?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	err := auth.New(&auth.Config{}).Init()
	if err != nil {
		logx.NewTraceLogger(context.Background()).Info().Msg(fmt.Sprintf("init err:%v", err))
	}
	mysql.InitAssetsMysql(env.DbDSN, 10)

	v11.SetUp()
	router := rest.NewGinServer()
	router.Use(httprouter.Recovery())
	router.Use(middleware.Cors())
	v1 := router.Group("/v1")
	{
		var (
			objective api.ObjectiveController
		)
		v1.POST("login", objective.Login)
		v1.Use(middleware.AuthenticatedHandlev2())

		v1.GET("test", objective.First)
		v1.GET("gongzhu", objective.GongZhu)
		v1.GET("redis", objective.TestRedis)
		v1.GET("mysql", objective.Mysql)
		v1.GET("rpc", objective.Rpc)
	}

	srv := &http.Server{
		Addr:         ":8086",
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
func setupGRPCServer() {
	go func() {
		if err := grpcbase.RegisterServer(grpcsvc.NewService()); err != nil {
			log.Fatal().Err(err).Msg("listen 8084 失败")
		}
	}()
}

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
