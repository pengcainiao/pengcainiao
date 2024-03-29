package internal

import (
	"context"
	middlewares "github.com/pengcainiao/pengcainiao/usercenter/internal/middleware"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/api"

	"github.com/gin-contrib/cors"

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
// @title 用户中心API文档
// @version 1.0
// @description 用户中心API文档
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

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(
		httprouter.Recovery(),
		middlewares.TokenVerify(),
	)

	v1 := router.Group("/v1")
	{
		var (
			user    api.UserController
			account = api.NewAccountController()
		)

		v1.GET("test", user.First)
		v1.POST("/user/register", user.RegisterUser) // 注册
		v1.POST("/user/login", user.UserLogin) // 登陆

		v1.GET("/auth/verify", account.VerifyAccessibleHandler) // 校验用户信息

	}

	srv := &http.Server{
		Addr:         ":8081",
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
