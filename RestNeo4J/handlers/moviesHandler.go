package handlers

import (
	"Rest/data"
	"Rest/domain"
	"encoding/json"

	// "context"
	"log"
	"net/http"
	// "strconv"
	// "github.com/gorilla/mux"
)

type KeyProduct struct{}

type UserHandler struct {
	logger *log.Logger
	// NoSQL: injecting movie repository
	repo *data.UserRepository
}

// Injecting the logger makes this code much more testable.
func NewUserHandler(l *log.Logger, r *data.UserRepository) *UserHandler {
	return &UserHandler{l, r}
}

// func (m *UserHandler) GetAllMovies(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	limit, err := strconv.Atoi(vars["limit"])
// 	if err != nil {
// 		m.logger.Printf("Expected integer, got: %d", limit)
// 		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
// 		return
// 	}

// 	movies, err := m.repo.GetAllNodesWithMovieLabel(limit)
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 	}

// 	if movies == nil {
// 		return
// 	}

// 	err = movies.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		m.logger.Fatal("Unable to convert to json :", err)
// 		return
// 	}
// }

// func (m *MoviesHandler) GetAllMoviesWithCast(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	limit, err := strconv.Atoi(vars["limit"])
// 	if err != nil {
// 		m.logger.Printf("Expected integer, got: %d", limit)
// 		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
// 		return
// 	}

// 	movies, err := m.repo.GetAllMoviesWithCast(limit)
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 	}

// 	if movies == nil {
// 		return
// 	}

// 	err = movies.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		m.logger.Fatal("Unable to convert to json :", err)
// 		return
// 	}
// }

// func (m *MoviesHandler) GetAllMoviesByTitle(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	title := vars["title"]

// 	movies, err := m.repo.GetAllNodesWithMovieLabelAndGivenTitle(title)
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 	}

// 	if movies == nil {
// 		return
// 	}

// 	err = movies.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		m.logger.Fatal("Unable to convert to json :", err)
// 		return
// 	}
// }

func (m *UserHandler) CreateUser(rw http.ResponseWriter, h *http.Request) {
	var person domain.User
	err := json.NewDecoder(h.Body).Decode(&person)
	if err != nil {
		m.logger.Print("Can't decode request body: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = m.repo.WriteUser(&person)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

// func (m *MoviesHandler) GetActorRole(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	role := vars["role"]

// 	roles, err := m.repo.GetPersonWhoPlayedRole(role)
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 	}

// 	if roles == nil {
// 		return
// 	}

// 	err = roles.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		m.logger.Fatal("Unable to convert to json :", err)
// 		return
// 	}
// }

// func (m *MoviesHandler) CreateMovie(rw http.ResponseWriter, h *http.Request) {
// 	movie := h.Context().Value(KeyProduct{}).(*data.Movie)
// 	err := m.repo.WriteMovie(movie)
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	rw.WriteHeader(http.StatusCreated)
// }

// func (m *MoviesHandler) GetPersonWhoActedAndProducedMovie(rw http.ResponseWriter, h *http.Request) {
// 	roles, err := m.repo.GetPersonWhoActedAndProducedMovie()
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 	}

// 	if roles == nil {
// 		return
// 	}

// 	err = roles.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		m.logger.Fatal("Unable to convert to json :", err)
// 		return
// 	}
// }

// func (m *MoviesHandler) GetKeanuMovies(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	limit, err := strconv.Atoi(vars["limit"])
// 	if err != nil {
// 		m.logger.Printf("Expected integer, got: %d", limit)
// 		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
// 		return
// 	}

// 	movies, err := m.repo.GetKeanuMovies(limit)
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 	}

// 	if movies == nil {
// 		return
// 	}

// 	err = movies.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		m.logger.Fatal("Unable to convert to json :", err)
// 		return
// 	}
// }

// func (m *MoviesHandler) GetActorsWithMostMovies(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	limit, err := strconv.Atoi(vars["limit"])
// 	if err != nil {
// 		m.logger.Printf("Expected integer, got: %d", limit)
// 		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
// 		return
// 	}

// 	actorMovies, err := m.repo.GetActorsWithMostMovies(limit)
// 	if err != nil {
// 		m.logger.Print("Database exception: ", err)
// 	}

// 	if actorMovies == nil {
// 		return
// 	}

// 	err = actorMovies.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		m.logger.Fatal("Unable to convert to json :", err)
// 		return
// 	}
// }

// func (m *MoviesHandler) MiddlewarePersonDeserialization(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
// 		person := &data.Person{}
// 		err := person.FromJSON(h.Body)
// 		if err != nil {
// 			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
// 			m.logger.Fatal(err)
// 			return
// 		}
// 		ctx := context.WithValue(h.Context(), KeyProduct{}, person)
// 		h = h.WithContext(ctx)
// 		next.ServeHTTP(rw, h)
// 	})
// }

// func (m *MoviesHandler) MiddlewareMovieDeserialization(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
// 		movie := &data.Movie{}
// 		err := movie.FromJSON(h.Body)
// 		if err != nil {
// 			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
// 			m.logger.Fatal(err)
// 			return
// 		}
// 		ctx := context.WithValue(h.Context(), KeyProduct{}, movie)
// 		h = h.WithContext(ctx)
// 		next.ServeHTTP(rw, h)
// 	})
// }

func (m *UserHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		m.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
