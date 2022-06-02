package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mystman.com/animated-pancake/internal/data"
	"mystman.com/animated-pancake/internal/handler"
	"mystman.com/animated-pancake/internal/service"
)

// LDFLAGS build tags:
var buildName string = "dev_build"
var buildVersion string = "default_version"
var buildDate string = "-"

// Application settings:
// TODO - potentially can take this from args using flag package
var dbFilePath = "/usr/share/pancake-data/anim-pancake.db"
var port = ":6543"

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

	// Channels for OS signal handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	// Start listening for requests
	go func() {
		repo := data.NewRepo(dbFilePath)
		svc := service.NewService(repo)
		hnd := handler.NewHandler(svc)

		mux := http.NewServeMux()
		mux.HandleFunc("/v1/", hnd.HandleRoot)
		mux.HandleFunc(fmt.Sprintf("/v1/%s", data.NetworkType), hnd.HandleNetwork)

		// For liveness probe
		mux.HandleFunc("/v1/liveness", handler.HandleReadiness)

		server := &http.Server{
			Addr:    port,
			Handler: mux,
		}

		log.Printf("Starting service on port %v", port)
		err := server.ListenAndServe()
		log.Fatal("Fatal error:", err)
	}()

	// Waiting for signal
	go func() {
		for {
			signal := <-shutdown
			log.Printf("Signal received: %v", signal)
			done <- true
		}
	}()

	// Waiting for shutdown signal
	<-done
	log.Printf("Shutting down service")

	return nil
}
