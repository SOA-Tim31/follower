package handlers

import (
	"Rest/data"
	"Rest/domain"
	"context"
	"encoding/json"

	// "context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "strconv"
	// "github.com/gorilla/mux"
)

type KeyProduct struct{}

type UserHandler struct {
	logger *log.Logger

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

func (m *UserHandler) FollowUser(rw http.ResponseWriter, h *http.Request) {

	var users []domain.User
	err := json.NewDecoder(h.Body).Decode(&users)

	if err != nil {
		m.logger.Print("Can't decode request body: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) < 2 {
		http.Error(rw, "Potrebno je poslati barem dve osobe", http.StatusBadRequest)
		return
	}

	err = m.repo.FollowUser(&users[0], &users[1])
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (f *UserHandler) GetRecommendations(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["userId"]

	// Assuming GetRecommendations should return a slice of domain.User
	users, err := f.repo.GetRecommendations(id)
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		f.logger.Print("Database exception: ", err)
		return
	}

	if users == nil || len(users) == 0 {
		jsonData, err := json.Marshal([]domain.User{})
		if err != nil {
			http.Error(rw, "Error marshaling json", http.StatusInternalServerError)
			f.logger.Print("Error marshaling json: ", err)
			return
		}
		rw.Write(jsonData)
		return
	}

	if err := json.NewEncoder(rw).Encode(users); err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		f.logger.Fatal("Unable to convert to json :", err)
	}
}

func (m *UserHandler) MiddlewareFollowingDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		newFollower := &domain.UserFollower{}
		err := newFollower.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			m.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, newFollower)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}

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
