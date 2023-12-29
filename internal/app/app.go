// Package app configures and runs application.
package app

import (
	"fmt"
	firestorerepo "lemon_be/internal/usecase/repo/firestoreRepo"
	"lemon_be/internal/usecase/repo/redisrepo"
	"lemon_be/internal/util/jwt"
	"lemon_be/pkg/redispkg"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"lemon_be/config"
	v1 "lemon_be/internal/controller/http/v1"
	"lemon_be/internal/usecase"
	"lemon_be/pkg/firestore"
	"lemon_be/pkg/httpserver"
	"lemon_be/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// db
	firestoreDb, err := firestore.NewFirestore(cfg.Firestore.ServiceAccLocation, cfg.Firestore.ServiceAccLocation)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - firestoreDb - firestore.NewFirestore: %w", err))
	}
	redis, err := redispkg.NewRedis(cfg.Redis.Address, cfg.Redis.Password)

	//repo
	userRepo := firestorerepo.NewUserRepo(firestoreDb)
	sessionRepo := firestorerepo.NewSessionRepo(firestoreDb)
	geoRedisRepo := redisrepo.NewGeoRedisRepo(redis)
	userRedisRepo := redisrepo.NewUserRedisRepo(redis)
	helpRepo := firestorerepo.NewHelpRepo(firestoreDb)
	// jwt
	jwtTokenMaker, err := jwt.NewJWTMaker("VBKNhRGFYZWGtbQ8hQ6ABQn1oNbYkHTu/fj/cUUO9p8=")

	// usecase
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtTokenMaker, sessionRepo)
	caregiverUseCase := usecase.NewCaregiverUseCase(geoRedisRepo, userRedisRepo, helpRepo)
	hub := usecase.NewHub(redis, geoRedisRepo)
	go hub.Run()
	websocketUSecase := usecase.NewWebsocketUseCase(hub, userRepo, geoRedisRepo)
	userUsecase := usecase.NewUserUseCase(userRedisRepo)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, authUseCase, websocketUSecase, caregiverUseCase, userUsecase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	firestoreDb.Client.Close()
	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
