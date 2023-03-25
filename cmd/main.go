package main

import (
	"context"
	"github.com/kevinicky/go-guest-book/delivery"
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

	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalln("error while connecting to database:", err)
	}

	guestBookRepository := newGuestBookRepository(db)
	guestBookUseCase := newGuestBookUseCase(guestBookRepository)
	guestBookAdapter := newGuestBookAdapter(guestBookUseCase)

	mux := http.NewServeMux()
	h := delivery.HTTPHandler{}
	h.NewRest(mux, guestBookAdapter)

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
		Handler:      mux,
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
