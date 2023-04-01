package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kevinicky/go-guest-book/delivery"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/repository/database"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"
)

func main() {
	configPath, err := filepath.Abs(path.Join("config"))
	if err != nil {
		log.Fatalln("error while opening config folder:", err)
	}

	viper.SetConfigFile(configPath + "./config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		log.Fatalln("error while reading config.yaml:", err)
	}

	pgDB, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalln("error while connecting to database postgresql:", err)
	}

	dbRedis, _, err := database.NewRedisDB()
	if err != nil {
		log.Fatalln("error while connecting to database redis:", err)
	}

	healthRepository := newHealthRepository(pgDB, dbRedis)
	healthUseCase := newHealthUseCase(healthRepository)
	healthAdapter := newHealthAdapter(healthUseCase)

	userRepository := newUserRepository(pgDB)
	userUseCase := newUserUseCase(userRepository)
	userAdapter := newUserAdapter(userUseCase)

	visitRepository := newVisitRepository(pgDB)
	visitUseCase := newVisitUseCase(visitRepository, userUseCase)
	visitAdapter := newVisitAdapter(visitUseCase, userAdapter)

	threadRepository := newThreadRepository(pgDB)
	threadUseCase := newThreadUseCase(threadRepository, visitUseCase, userUseCase)
	threadAdapter := newThreadAdapter(threadUseCase, visitAdapter, userAdapter)

	jwtSecretKey := []byte(viper.GetString("jwt.secret"))
	jwtExpired := viper.GetDuration("jwt.expired")
	authUseCase := newAuthUseCase(userUseCase, entity.JwtClaims{
		Expired:   jwtExpired * time.Minute,
		SecretKey: jwtSecretKey,
	})
	authAdapter := newAuthAdapter(authUseCase)

	r := mux.NewRouter()
	h := delivery.HTTPHandler{}
	h.NewRest(r, healthAdapter, userAdapter, visitAdapter, threadAdapter, authAdapter)

	appName := viper.GetString("app.name")
	appServer := viper.GetString("app.server")
	appPort := viper.GetString("app.port")
	appReadTO := viper.GetDuration("app.timeout.read")
	appWriteTO := viper.GetDuration("app.timeout.write")
	appIdleTO := viper.GetDuration("app.timeout.idle")

	server := &http.Server{
		Addr:         appServer + ":" + appPort,
		ReadTimeout:  appReadTO * time.Second,
		WriteTimeout: appWriteTO * time.Second,
		IdleTimeout:  appIdleTO * time.Second,
		Handler:      r,
	}

	go func() {
		log.Println("running " + appName + " on " + server.Addr)
		if errServe := server.ListenAndServe(); errServe != nil {
			log.Println(errServe.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if errShutDown := server.Shutdown(context.Background()); errShutDown != nil {
		log.Println(errShutDown.Error())
	}

	log.Println("shutting down...")
	os.Exit(0)
}
