package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Simcha-b/Podcast-Hub/models"
	"github.com/Simcha-b/Podcast-Hub/utils"
	"github.com/mmcdole/gofeed"
)

var Logger = utils.NewLogger("info")

func parseRSSFeed(url string) ([]models.Episode, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Fatalf("Error parsing RSS feed: %v", err)
	}
	Logger.Info(fmt.Sprintf("Successfully parsed RSS feed: %s", feed.Title))
	var episodes []models.Episode
	for _, item := range feed.Items {
		episode := models.Episode{
			ID:          item.GUID,
			PodcastID:   time.Now().Format("20060102150405"), // Placeholder for PodcastID, should be replaced with actual logic
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			AudioURL:    item.Enclosures[0].URL, // Assuming the first enclosure is the audio file
			// Duration:   item.Duration,
			// PublishedAt: item.PublishedParsed,
			// FileSize:    item.Enclosures[0].Length, // Assuming the first enclosure is the audio file
			CreatedAt: time.Now(), // Set to current time, adjust as needed
		}
		episodes = append(episodes, episode)
	}
	Logger.Info(fmt.Sprintf("Parsed %d episodes from the feed", len(episodes)))
	fmt.Println("Parsed episodes:", episodes[0]) //TODO: Remove this line in production
	return episodes, nil
}
