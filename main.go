package main

import (
	"github.com/rajatgpt1521/cachingSystem/service/models"
	"github.com/rajatgpt1521/cachingSystem/service/pkg/cache_handler"
	"github.com/rajatgpt1521/cachingSystem/service/pkg/database"
	server2 "github.com/rajatgpt1521/cachingSystem/service/server"
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	database.Initialize()
	models.AutoMigrateSQL()
	cache_handler.Initialize()

}
func main() {

	log.Info().Msg("Start Cache")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	r := mux.NewRouter()
	server := http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
	}
	r.HandleFunc("/", Routes).Methods("GET")
	r.HandleFunc("/view/page/{limit}", server2.ReadCachePagination).Methods("GET")
	r.HandleFunc("/insert/{data}", server2.PutData).Methods("PUT")
	r.HandleFunc("/notify/{msg}", server2.Reload).Methods("PUT")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Server closed unexpectedly.")
		}
	}()
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	// closing DB connection
	defer database.Close()
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info().Msg("shutting down")
	os.Exit(0)

}

func Routes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"Put" : "/insert/{data}","Read" : "/view/page/{pageid}"}`)
}
