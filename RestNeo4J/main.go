package main

import (
	"Rest/data"
	"Rest/handlers"
	follower "Rest/proto"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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


	postFollowBranch := router.Methods(http.MethodPost).Subrouter()
    postFollowBranch.HandleFunc("/follower", userHandler.FollowUser)

	//GRPC
	listener, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)
	logger.Println("Server listening on port 89")

	userHandlergRPC := handlers.NewgRPCUserHandler(store)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	follower.RegisterFollowerServiceServer(grpcServer, userHandlergRPC)

	//Distribute all the connections to goroutines
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	logger.Println("Server stopped")
}
