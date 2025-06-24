package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Simcha-b/Podcast-Hub/services"
	"github.com/gorilla/mux"
)

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
	// החזרת כל הפרקים עם אפשרות לסינון
}

func Search(w http.ResponseWriter, r *http.Request) {
	// חיפוש פודקאסטים ופרקים
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	// החזרת סטטיסטיקות מערכת
}

func AddFeed(w http.ResponseWriter, r *http.Request) {
	// הוספת מקור RSS חדש
}

func DeleteFeed(w http.ResponseWriter, r *http.Request) {
	// מחיקת מקור RSS לפי כתובת
}

func GetDailyReport(w http.ResponseWriter, r *http.Request) {
	// הפקת דוח יומי
}
