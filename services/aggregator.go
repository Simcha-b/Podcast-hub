package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Simcha-b/Podcast-Hub/models"
)

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

// AggregateAllFeeds מריץ את כל הפידים אחד אחד (או מקבילי בהמשך) ועושה parsing ושמירה
func AggregateAllFeeds(storage *FileStorage, feedSources []models.Feed) error {

	for _, feed := range feedSources {
		if err := ProcessSingleFeed(storage, feed); err != nil {
			Logger.Error(fmt.Sprintf("Error processing feed %s: %v", feed.URL, err))
			if err := UpdateFeedStatus(feed, false); err != nil {
				Logger.Error(fmt.Sprintf("Failed to update feed status for %s: %v", feed.URL, err))
			}
			continue
		}
		if err := UpdateFeedStatus(feed, true); err != nil {
			Logger.Error(fmt.Sprintf("Failed to update feed status for %s: %v", feed.URL, err))
		}
	}
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

	// Save podcast and episodes to storage
	if err := storage.SavePodcast(podcast); err != nil {
		Logger.Error(fmt.Sprintf("Failed to save podcast %s: %v", podcast.ID, err))
		return err
	}
	for _, episode := range episodes {
		if err := storage.SaveEpisode(&episode); err != nil {
			Logger.Error(fmt.Sprintf("Failed to save episode %s for podcast %s: %v", episode.ID, podcast.ID, err))
			return err
		}
	}
	Logger.Info(fmt.Sprintf("Successfully processed feed %s with %d episodes", feed.URL, len(episodes)))
	// Update feed status
	return nil
}

func RunAggregator() {
	// Load feed sources from JSON file
	feedSources, err := LoadFeedSources("data/feeds.json")
	if err != nil {
		Logger.Error(fmt.Sprintf("Failed to load feed sources: %v", err))
		return
	}

	// Initialize storage (assuming FileStorage is defined elsewhere)
	storage := NewFileStorage("data")

	// Aggregate all feeds
	if err := AggregateAllFeeds(storage, feedSources); err != nil {
		Logger.Error(fmt.Sprintf("Error aggregating feeds: %v", err))
	}
}
