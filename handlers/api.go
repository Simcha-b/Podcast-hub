package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Simcha-b/Podcast-Hub/config"
	"github.com/Simcha-b/Podcast-Hub/services"
	"github.com/Simcha-b/Podcast-Hub/utils"

	"github.com/gorilla/mux"
)

var cfg = config.LoadConfig()

var Logger = utils.NewLogger(cfg.LOG_LEVEL)

var storage = services.NewFileStorage(cfg.DATA_DIR)

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
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "Podcast not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch podcast", http.StatusInternalServerError)
		return
	}
	// if podcast == nil {
	// 	http.Error(w, "Podcast not found", http.StatusNotFound)
	// 	return
	// }
	writeJSON(w, http.StatusOK, podcast)
}

func GetPodcastEpisodes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podcastId := vars["id"]
	episodes, err := storage.LoadEpisodes(podcastId)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			Logger.Info("No podcast found with ID: " + podcastId)
			http.Error(w, "No episodes found for this podcast", http.StatusNotFound)
			return
		}
		Logger.Error("Failed to fetch episodes for podcast " + podcastId + ": " + err.Error())
		http.Error(w, "Failed to fetch episodes", http.StatusInternalServerError)
		return
	}
	if len(episodes) == 0 {
		http.Error(w, "No episodes found for this podcast", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, episodes)
}

func GetEpisodeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["podcastId"] == "" || vars["episodeId"] == "" {
		http.Error(w, "Podcast ID and Episode ID are required", http.StatusBadRequest)
		return
	}
	podcastId := vars["podcastId"]
	episodeId := vars["episodeId"]

	episode, err := storage.LoadEpisodeByID(podcastId, episodeId)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			Logger.Info("Episode not found with ID: " + episodeId)
			http.Error(w, "Episode not found", http.StatusNotFound)
			return
		}
		Logger.Error("Failed to fetch episode: " + err.Error())
		http.Error(w, "Failed to fetch episode", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, episode)
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
		if errors.Is(err, services.ErrNotFound) {
			Logger.Info("No podcasts found for query: " + query)
			http.Error(w, "No podcasts found for the given query", http.StatusNotFound)
			return
		}
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
	services.ProcessSingleFeed(storage, *feed)
	Logger.Info("Successfully added new feed: " + feed.URL)
	writeJSON(w, http.StatusCreated, feed)
}

func DeleteFeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Feed ID is required", http.StatusBadRequest)
		return
	}
	if err := services.DeleteFeedFromSources(url); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			Logger.Info("Feed not found for URL: " + url)
			http.Error(w, "Feed not found", http.StatusNotFound)
			return
		}
		Logger.Error("Failed to delete feed: " + err.Error())
		http.Error(w, "Failed to delete feed", http.StatusInternalServerError)
		return
	}
	Logger.Info("Successfully deleted feed with URL: " + url)
	w.WriteHeader(http.StatusNoContent)
}

func GetLastEpisodes(w http.ResponseWriter, r *http.Request) {
	episods, err := storage.LoadLastEpisodes()
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			Logger.Info("No episodes found")
			http.Error(w, "No episodes found", http.StatusNotFound)
			return
		}
		Logger.Error("Failed to fetch last episodes: " + err.Error())
		http.Error(w, "Failed to fetch last episodes", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, episods)
}

func GetDownloadLink(w http.ResponseWriter, r *http.Request) {
	podcastId := r.URL.Query().Get("podcastId")
	episodeId := r.URL.Query().Get("episodeId")

	episode, err := storage.LoadEpisodeByID(podcastId, episodeId)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			Logger.Info("Episode not found with ID: " + episodeId)
			http.Error(w, "Episode not found", http.StatusNotFound)
			return
		}
		Logger.Error("Failed to fetch episode: " + err.Error())
		http.Error(w, "Failed to fetch episode", http.StatusInternalServerError)
		return
	}
	if episode.AudioURL == "" {
		Logger.Info("No audio URL found for episode ID: " + episodeId)
		http.Error(w, "No audio URL found for this episode", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"download_link": episode.AudioURL})
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	podcasts, err := storage.LoadAllPodcasts()
	if err != nil {
		http.Error(w, "Failed to fetch podcasts", http.StatusInternalServerError)
		Logger.Error("Failed to fetch podcasts")
		return
	}

	totalEpisodes := 0
	for _, podcast := range podcasts {
		totalEpisodes += podcast.NumOfEpisodes
	}

	stats := map[string]int{
		"podcast_count": len(podcasts),
		"episode_count": totalEpisodes,
	}
	Logger.Info(fmt.Sprintf("Stats result: %v", stats))
	writeJSON(w, http.StatusOK, stats)

}
