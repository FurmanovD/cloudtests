package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/FurmanovD/cloudtests/internal/app/userservice"
	"github.com/FurmanovD/cloudtests/internal/app/webserver"
	"github.com/FurmanovD/cloudtests/internal/pkg/config"
	"github.com/FurmanovD/cloudtests/internal/pkg/errortranslator"
	"github.com/FurmanovD/cloudtests/internal/pkg/repository"
)

func main() {

	// Config
	appConf, dbConf, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(appConf.Loglevel)
	log.Info("service started")
	defer log.Info("service exited")

	// Repository
	log.Info("creating a DB repository")
	repository, err := repository.New(dbConf.Address, dbConf.Password)
	if err != nil {
		log.Fatalf("repository creation error: %v", err)
	}
	log.Info("DB repository has been created")

	// Service
	userService := userservice.New(
		repository,
		errortranslator.NewRepositoryErrorTranslator(),
	)

	// manage a process interruption as well
	finishChannel := make(chan error)
	defer close(finishChannel)

	// SIGINT, SIGTERM handling
	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		finishChannel <- fmt.Errorf("%s", <-c)
	}()

	// start listening
	// HTTP server
	httpHandler := webserver.NewHTTPServer(userService)
	go func() {
		server := &http.Server{
			Addr:        appConf.ServiceAddress,
			Handler:     httpHandler,
			ReadTimeout: time.Duration(appConf.HTTPTimeoutSec) * time.Second,
		}
		log.Infof("HTTP server started listening %s", appConf.ServiceAddress)
		finishChannel <- server.ListenAndServe()
	}()

	err = <-finishChannel
	fmt.Println(err.Error())
}
