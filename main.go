package main

import (
	// "fmt"
	"net/http"

	"github.com/Simcha-b/Podcast-Hub/handlers"
	"github.com/Simcha-b/Podcast-Hub/services"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
)

func main() {

	// fmt.Println("welcome to the Podcast-Hub!!")
	myFigure := figure.NewColorFigure("Podcast-Hub!!", "", "red", true)
	myFigure.Print()

	go services.RunAggregator()

	r := mux.NewRouter()
	r.HandleFunc("/api/podcasts", handlers.GetPodcasts).Methods("GET")
	r.HandleFunc("/api/podcasts/{id}", handlers.GetPodcastByID).Methods("GET")
	r.HandleFunc("/api/podcasts/{id}/episodes", handlers.GetPodcastEpisodes).Methods("GET")
	r.HandleFunc("/api/podcast/{podcastId}/{episodeId}", handlers.GetEpisodeByID).Methods("GET")
	r.HandleFunc("/api/lastEpisodes", handlers.GetLastEpisodes).Methods("GET")

	r.HandleFunc("/api/search", handlers.Search).Methods("GET")
	r.HandleFunc("/api/download", handlers.GetDownloadLink).Methods("GET")

	r.HandleFunc("/api/feeds", handlers.AddFeed).Methods("GET")
	r.HandleFunc("/api/feeds", handlers.DeleteFeed).Methods("DELETE")
	r.HandleFunc("/api/stats", handlers.GetStats).Methods("GET")

	// Starting the server
	handlers.Logger.Info("Starting server on port 8080")
	http.ListenAndServe(":8080", r)
}
