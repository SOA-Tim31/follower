package main

import (
	"Rest/data"
	"Rest/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// "github.com/gorilla/mux"
)

func main() {
	//Reading from environment, if not set we will default it to 8080.
	//This allows flexibility in different environments (for eg. when running multiple docker api's and want to override the default port)
	

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[movie-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[movie-store] ", log.LstdFlags)

	// NoSQL: Initialize Movie Repository store
	store, err := data.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseDriverConnection(timeoutContext)
	store.CheckConnection()

	userHandler := handlers.NewUserHandler(logger, store)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()

	router.Use(userHandler.MiddlewareContentTypeSet)

	// getMovieByTitle := router.Methods(http.MethodGet).Subrouter()
	// getMovieByTitle.HandleFunc("/movies/title/{title}", moviesHandler.GetAllMoviesByTitle)

	// getAllMoviesWithCast := router.Methods(http.MethodGet).Subrouter()
	// getAllMoviesWithCast.HandleFunc("/movies/cast/{limit}", moviesHandler.GetAllMoviesWithCast)

	// getAllMovies := router.Methods(http.MethodGet).Subrouter()
	// getAllMovies.HandleFunc("/movies/{limit}", moviesHandler.GetAllMovies)

	postUserNode := router.Methods(http.MethodPost).Subrouter()
	postUserNode.HandleFunc("/user", userHandler.CreateUser)
	postUserNode.Use(userHandler.MiddlewareContentTypeSet)

	

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":89",
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Println("Server listening on port 89")
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
