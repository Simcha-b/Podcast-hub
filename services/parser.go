package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	. "github.com/Simcha-b/Podcast-Hub/models"
)

type Parser interface {
	ParseFeed(url string) (*Podcast, []*Episode, error)
}

type RSSParser struct{}

func (p *RSSParser) ParseFeed(url string) (*Podcast, []*Episode, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rss Feed
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	podcast := &Podcast{
		ID:          rss.Channel.ID,
		Title:       rss.Channel.Title,
		Description: rss.Channel.Description,
		ImageURL:    rss.Channel.Image.URL,
		FeedURL:     url,
	}

	var episodes []*Episode
	for _, item := range rss.Channel.Items {
		pubDate, _ := time.Parse(time.RFC1123Z, item.PubDate)
		episodes = append(episodes, &Episode{
			ID:          item.GUID,
			Name:        item.Title,
			Description: item.Description,
			PubDate:     pubDate,
			URL:         item.Enclosure.URL,
			Duration:    item.Enclosure.Length,
			Type:        item.Enclosure.Type,
			ImageURL:    item.Image.URL,
		})
	}

	return podcast, episodes, nil
}
