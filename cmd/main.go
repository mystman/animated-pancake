package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"mystman.com/animated-pancake/internal/data"
	"mystman.com/animated-pancake/internal/service"
)

// LDFLAGS build tags:
var buildName string = "dev_build"
var buildVersion string = "default_version"
var buildDate string = "-"

// Application settings:
// TODO - potentially can take this from args using flag package
var dbFilePath = "/usr/share/pancake-data/anim-pancake.db"

func main() {

	log.Printf("Staring main %v with params: %v", os.Args[0], os.Args[1:])
	log.Printf("%v :: verion %v [ built: %v]\n", buildName, buildVersion, buildDate)

	if err := Run(); err != nil {
		log.Fatalf("Fatality: %v", err)
	}

	log.Printf("Stopping %v", os.Args[0])

}

// Run - instantiation and startup
func Run() error {
	log.Printf("Staring service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	// starting the service
	repo := data.NewRepository(dbFilePath)
	svc := service.NewService(repo)

	// service
	go func(svc *service.Service) {
		for {
			signal := <-shutdown
			log.Printf("Signal received: %v", signal)

			done <- true
		}
	}(svc)
	log.Printf("Service has started")

	// Waiting for shutdown signal
	<-done
	log.Printf("Shutting down service")

	return nil
}
