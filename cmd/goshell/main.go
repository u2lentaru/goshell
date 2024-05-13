package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goshell/internal/pgclient"
	"goshell/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	// cfg := config.MustLoad()
	//$ENV:CONFIG_PATH="E:\workgo\goshell\config\local.yaml"
	//$ENV:DB_HOST="localhost"

	// TODO: logs

	url := "postgres://postgres:postgres@" + os.Getenv("DB_HOST") + ":5432/postgres"

	dbpool, err := pgclient.GetDb(context.Background(), url)
	defer dbpool.Close()

	if err != nil {
		log.Fatal(err)
	}

	route := mux.NewRouter()

	routes.AddRoutes(route)

	//GS start
	srv := &http.Server{
		Addr:    ":8080",
		Handler: route,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started at http://localhost:8080/")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		log.Println("Sleep on")
		time.Sleep(time.Second * 1)
		log.Println("Sleep off")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
