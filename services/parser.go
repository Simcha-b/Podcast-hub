package services

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Simcha-b/Podcast-Hub/models"
	"github.com/Simcha-b/Podcast-Hub/utils"
	"github.com/mmcdole/gofeed"
)

// Logger instance for logging within the parser service

var Logger = utils.NewLogger(cfg.LOG_LEVEL)

func generatePodcastID(feedURL string) string {
	h := sha1.New()
	h.Write([]byte(feedURL))
	return hex.EncodeToString(h.Sum(nil))[:12] // מזהה קצר
}

func ParseNewFeedSources(url string) (*models.Feed, error) {
	// Create a new Feed instance with the provided URL
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(url)
	if err != nil {
		// Log and return error if RSS feed parsing fails
		Logger.Error(fmt.Sprintf("Failed to parse RSS feed from URL %s: %v", url, err))
		return &models.Feed{}, err
	}
	// Create a new Feed struct from the parsed feed data
	newfeed := &models.Feed{
		URL:         url,
		Name:        feed.Title, // Optional: can be set later
		Active:      true,
		LastFetched: time.Now(),
		ErrorCount:  0,
	}
	log.Printf(feed.Updated)
	Logger.Info(fmt.Sprintf("Successfully parsed RSS feed: %s", feed.Title))
	return newfeed, nil
}

// parseRSSFeed parses an RSS feed from the given URL and returns a Podcast and its Episodes
func parseRSSFeed(url string) (*models.Podcast, []models.Episode, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(url)
	if err != nil {
		// Log and return error if RSS feed parsing fails
		Logger.Error(fmt.Sprintf("Failed to parse RSS feed from URL %s: %v", url, err))
		return nil, nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}
	Logger.Info(fmt.Sprintf("Successfully parsed RSS feed: %s", feed.Title))
	// Build Podcast struct from feed data
	podcastID := generatePodcastID(url)
	podcast := &models.Podcast{
		ID:          podcastID,
		Title:       feed.Title,
		Description: feed.Description,
		Author:      feed.Author.Name,
		ImageURL:    feed.Image.URL,
		FeedURL:     url,
		Category:    "",
		Language:    feed.Language,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	var episodes []models.Episode
	// Iterate over feed items to build Episode structs
	for _, item := range feed.Items {
		episode := models.Episode{
			ID:          item.GUID,
			PodcastID:   podcastID,
			Title:       item.Title,
			Description: item.Description,
			CreatedAt:   time.Now(),
		}

		// Set PublishedAt if available in the feed item
		if item.PublishedParsed != nil {
			episode.PublishedAt = *item.PublishedParsed
		}

		// Extract audio URL and file size from enclosure if present
		if len(item.Enclosures) > 0 {
			episode.AudioURL = item.Enclosures[0].URL
			// Convert file size from string to int64
			episode.FileSize, _ = strconv.ParseInt(item.Enclosures[0].Length, 10, 64)
		}
		episodes = append(episodes, episode)
	}
	podcast.NumOfEpisodes = len(episodes)
	Logger.Info(fmt.Sprintf("Parsed %d episodes from the feed", len(episodes)))
	return podcast, episodes, nil
}
