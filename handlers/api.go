package handlers

import (
	"encoding/json"
	"net/http"
	
	"github.com/Simcha-b/Podcast-Hub/services"
	"github.com/Simcha-b/Podcast-Hub/utils"
	"github.com/gorilla/mux"
)

var Logger = utils.NewLogger("info")

var storage = services.NewFileStorage("data")

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "JSON encode failed", http.StatusInternalServerError)
	}
}

func GetPodcasts(w http.ResponseWriter, r *http.Request) {
	podcasts, err := storage.LoadAllPodcasts()
	if err != nil {
		http.Error(w, "Failed to fetch podcasts", http.StatusInternalServerError)
		return
	}
	if len(podcasts) == 0 {
		http.Error(w, "No podcasts found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, podcasts)
}

func GetPodcastByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podcastId := vars["id"]
	podcast, err := storage.LoadPodcastByID(podcastId)
	if err != nil {
		http.Error(w, "Failed to fetch podcast", http.StatusInternalServerError)
		return
	}
	if podcast == nil {
		http.Error(w, "Podcast not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, podcast)
}

func GetPodcastEpisodes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podcastId := vars["id"]
	episodes, err := storage.LoadEpisodes(podcastId)
	if err != nil {
		http.Error(w, "Failed to fetch episodes", http.StatusInternalServerError)
		return
	}
	if len(episodes) == 0 {
		http.Error(w, "No episodes found for this podcast", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, episodes)
}

func GetEpisodes(w http.ResponseWriter, r *http.Request) {

}

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	Logger.Info("Search query received: " + query)
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		Logger.Error("Search query is empty")
		return
	}

	podcasts, err := storage.SearchPodcasts(query)
	if err != nil {
		Logger.Error("Failed to search podcasts: " + err.Error())
		http.Error(w, "Failed to search podcasts", http.StatusInternalServerError)
		return
	}

	if len(podcasts) == 0 {
		Logger.Info("No podcasts found for query: " + query)
		http.Error(w, "No podcasts found for the given query", http.StatusNotFound)
		return
	}
	Logger.Info("Search completed successfully, found " + string(len(podcasts)) + " podcasts")
	writeJSON(w, http.StatusOK, podcasts)
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	// החזרת סטטיסטיקות מערכת
}

// הוספת מקור RSS חדש
func AddFeed(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("url")
	if query == "" {
		http.Error(w, "Query parameter 'url' is required", http.StatusBadRequest)
		return
	}
	feed, err := services.ParseNewFeedSources(query)
	if err != nil {
		Logger.Error("Failed to parse RSS feed: " + err.Error())
		http.Error(w, "Failed to parse RSS feed", http.StatusInternalServerError)
		return
	}
	if feed.URL == "" {
		Logger.Error("Parsed feed URL is empty")
		http.Error(w, "Parsed feed is empty", http.StatusBadRequest)
		return
	}
	if err := services.AddFeedToSources(*feed); err != nil {
		Logger.Error("Failed to save feed: " + err.Error())
		http.Error(w, "Failed to save feed", http.StatusInternalServerError)
		return
	}
	Logger.Info("Successfully added new feed: " + feed.URL)
	writeJSON(w, http.StatusCreated, feed)
}

func DeleteFeed(w http.ResponseWriter, r *http.Request) {
	// מחיקת מקור RSS לפי כתובת
}

func GetDailyReport(w http.ResponseWriter, r *http.Request) {
	// הפקת דוח יומי
}
