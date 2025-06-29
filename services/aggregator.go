package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Simcha-b/Podcast-Hub/config"
	"github.com/Simcha-b/Podcast-Hub/models"
)


var cfg = config.LoadConfig()

func LoadFeedSources(path string) ([]models.Feed, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		Logger.Error(fmt.Sprintf("Failed to read feed sources from %s: %v", path, err))
		return nil, err
	}

	var feeds []models.Feed
	if err := json.Unmarshal(data, &feeds); err != nil {
		return nil, err
	}

	return feeds, nil
}

func UpdateFeedStatus(source models.Feed, success bool) error {
	feeds, err := LoadFeedSources("data/feeds.json")
	if err != nil {
		return fmt.Errorf("failed to load feed sources: %w", err)
	}
	for i, feed := range feeds {
		if feed.URL == source.URL {
			if success {
				feeds[i].ErrorCount = 0
				feeds[i].LastFetched = source.LastFetched
				Logger.Info(fmt.Sprintf("Successfully updated feed %s status last fath: %s", feed.URL, feeds[i].LastFetched.Format("2006-01-02 15:04:05")))
			} else {
				feeds[i].ErrorCount++
				if feeds[i].ErrorCount > 5 {
					feeds[i].Active = false // Disable feed after 5 consecutive errors
				}
			}
			break
		}
	}
	data, err := json.MarshalIndent(feeds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated feeds: %w", err)
	}
	if err := os.WriteFile("data/feeds.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write updated feeds to file: %w", err)
	}

	return nil
}

func AddFeedToSources(feed models.Feed) error {
	feeds, err := LoadFeedSources("data/feeds.json")
	if err != nil {
		return fmt.Errorf("failed to load feed sources: %w", err)
	}
	// Check if the feed already exists
	for _, existingFeed := range feeds {
		if existingFeed.URL == feed.URL {
			Logger.Info(fmt.Sprintf("Feed %s already exists, skipping addition", feed.URL))
			return nil // Feed already exists, no need to add
		}
	}
	feeds = append(feeds, feed)
	data, err := json.MarshalIndent(feeds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated feeds: %w", err)
	}
	if err := os.WriteFile("data/feeds.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write updated feeds to file: %w", err)
	}
	Logger.Info(fmt.Sprintf("Successfully added feed %s to sources", feed.URL))
	return nil
}

func DeleteFeedFromSources(feedURL string) error {
	feeds, err := LoadFeedSources("data/feeds.json")
	if err != nil {
		return fmt.Errorf("failed to load feed sources: %w", err)
	}
	var updatedFeeds []models.Feed
	found := false
	for _, feed := range feeds {
		if feed.URL == feedURL {
			found = true
		} else {
			updatedFeeds = append(updatedFeeds, feed)
		}
	}
	if !found {
		Logger.Error(fmt.Sprintf("Feed %s not found, nothing to delete", feedURL))
		return ErrNotFound
	}
	data, err := json.MarshalIndent(updatedFeeds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated feeds: %w", err)
	}
	if err := os.WriteFile("data/feeds.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write updated feeds to file: %w", err)
	}
	Logger.Info(fmt.Sprintf("Successfully deleted feed %s from sources", feedURL))
	return nil
}

// TODO: Implement a worker pool to process feeds concurrently
func AggregateAllFeeds(storage *FileStorage, feedSources []models.Feed) error {
	var wg sync.WaitGroup

	for _, feed := range feedSources {
		wg.Add(1)
		go func(feed models.Feed) {
			defer wg.Done()
			err := ProcessSingleFeed(storage, feed)
			if err != nil {
				Logger.Error(fmt.Sprintf("Error processing feed %s: %v", feed.URL, err))
				if err := UpdateFeedStatus(feed, false); err != nil {
					Logger.Error(fmt.Sprintf("Failed to update feed status for %s: %v", feed.URL, err))
				}
				return
			}
			if err := UpdateFeedStatus(feed, true); err != nil {
				Logger.Error(fmt.Sprintf("Failed to update feed status for %s: %v", feed.URL, err))
			}
		}(feed)
	}
	wg.Wait()
	Logger.Info("All feeds processed successfully")
	return nil // Placeholder return, implement actual logic
}

func ProcessSingleFeed(storage *FileStorage, feed models.Feed) error {
	podcast, episodes, err := parseRSSFeed(feed.URL)
	if err != nil {
		Logger.Error(fmt.Sprintf("Failed to process feed %s: %v", feed.URL, err))
		return err
	}
	if podcast == nil || len(episodes) == 0 {
		return fmt.Errorf("no valid podcast or episodes found for feed %s", feed.URL)
	}

	if err := storage.SavePodcast(podcast); err != nil {
		Logger.Error(fmt.Sprintf("Failed to save podcast %s: %v", podcast.ID, err))
		return err
	}

	// load existing episodes for the podcast
	existingEpisodes, err := storage.LoadEpisodes(podcast.ID)
	if err != nil {
		Logger.Info(fmt.Sprintf("No existing episodes found for podcast %s (new podcast)", podcast.ID))
		existingEpisodes = []models.Episode{}
	}

	// Create a map of existing episodes for quick lookup
	existingMap := make(map[string]models.Episode)
	for _, ep := range existingEpisodes {
		existingMap[ep.ID] = ep
	}

	//save new episodes
	newCount := 0
	for _, episode := range episodes {
		if _, exists := existingMap[episode.ID]; !exists {
			if err := storage.SaveEpisode(&episode); err != nil {
				Logger.Error(fmt.Sprintf("Failed to save episode %s for podcast %s: %v", episode.ID, podcast.ID, err))
				return err
			}
			newCount++
		}
	}

	if newCount == 0 {
		Logger.Info(fmt.Sprintf("No new episodes for podcast %s", podcast.ID))
	} else {
		Logger.Info(fmt.Sprintf("Added %d new episodes for podcast %s", newCount, podcast.ID))
	}

	return nil
}

func IsPodcastOrEpisodesUpdated(storage *FileStorage, podcast *models.Podcast, episodes []models.Episode) (bool, error) {
	existingPodcast, err := storage.LoadPodcastByID(podcast.ID)
	if err != nil {
		return true, nil
	}
	if existingPodcast.UpdatedAt.Before(podcast.UpdatedAt) {
		return true, nil
	}
	existingEpisodes, err := storage.LoadEpisodes(podcast.ID)
	if err != nil {
		return true, nil
	}
	existingMap := make(map[string]models.Episode)
	for _, ep := range existingEpisodes {
		existingMap[ep.ID] = ep
	}
	for _, ep := range episodes {
		existing, ok := existingMap[ep.ID]
		if !ok || existing.PublishedAt.Before(ep.PublishedAt) {
			return true, nil
		}
	}
	return false, nil
}

func RunAggregator() {
	
	Logger.Info("Starting RSS Aggregator with ticker...")

	interval, err := strconv.Atoi(cfg.TIME_INTERVAL)
	if err != nil {
		Logger.Error(fmt.Sprintf("Invalid TIME_INTERVAL: %v", err))
		interval = 30 // Default to 30 minutes if parsing fails
	}
	// Ticker that runs every 30 minutes (can be changed as needed)
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	// Initialize storage
	storage := NewFileStorage(cfg.DATA_DIR)

	// Function to perform the aggregation logic
	doAggregation := func() {
		// Load feed sources from JSON file
		feedSources, err := LoadFeedSources(cfg.DATA_DIR + "/feeds.json")
		if err != nil {
			Logger.Error(fmt.Sprintf("Failed to load feed sources: %v", err))
			return
		}

		Logger.Info(fmt.Sprintf("Processing %d feeds", len(feedSources)))

		// Aggregate all feeds
		if err := AggregateAllFeeds(storage, feedSources); err != nil {
			Logger.Error(fmt.Sprintf("Error aggregating feeds: %v", err))
		} else {
			Logger.Info("Aggregation completed successfully")
		}
	}

	// הרצה ראשונית מיד
	Logger.Info("Running initial aggregation...")
	doAggregation()

	// לולאה שמקשיבה לטיקר
	for range ticker.C {
		Logger.Info("Running scheduled aggregation...")
		doAggregation()
	}
}
