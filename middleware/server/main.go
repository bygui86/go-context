package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(guidMiddleware)
	router.HandleFunc("/ishealthy", handleIsHealthy).Methods(http.MethodGet)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

// add value to context
func guidMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		r = r.WithContext(context.WithValue(r.Context(), "uuid", id))
		next.ServeHTTP(w, r)
	})
}

// retrieve value from context
func handleIsHealthy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	id := r.Context().Value("uuid")
	log.Printf("[%v] Returning 200 - Healthy", id)
	_, err := w.Write([]byte("Healthy"))
	if err != nil {
		panic(err)
	}
}
