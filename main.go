package main

import (
	"fmt"
	"net/http"

	"github.com/Simcha-b/Podcast-Hub/handlers"
	"github.com/Simcha-b/Podcast-Hub/services"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("welcome to the Podcast-Hub!!")
	services.RunAggregator()

	r := mux.NewRouter()
	r.HandleFunc("/api/podcasts", handlers.GetPodcasts).Methods("GET")
	r.HandleFunc("/api/podcasts/{id}", handlers.GetPodcastByID).Methods("GET")
	r.HandleFunc("/api/podcasts/{id}/episodes", handlers.GetPodcastEpisodes).Methods("GET")
	r.HandleFunc("/api/episodes", handlers.GetEpisodes).Methods("GET")
	r.HandleFunc("/api/search", handlers.Search).Methods("GET")
	r.HandleFunc("/api/stats", handlers.GetStats).Methods("GET")
	r.HandleFunc("/api/feeds", handlers.AddFeed).Methods("GET")
	r.HandleFunc("/api/feeds/{url}", handlers.DeleteFeed).Methods("DELETE")
	r.HandleFunc("/api/reports/daily", handlers.GetDailyReport).Methods("GET")

	// Starting the server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
