package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Simcha-b/Podcast-Hub/models"
)

type Storage interface {
	SavePodcast(podcast *models.Podcast) error
	LoadPodcast(id string) (*models.Podcast, error)
	SaveEpisode(episode *models.Episode) error
	LoadEpisode(podcastID, episodeID string) (*models.Episode, error)
}

type FileStorage struct {
	dataDir string
}

func NewFileStorage(dataDir string) *FileStorage {
	return &FileStorage{
		dataDir: dataDir,
	}
}

func (fs *FileStorage) SavePodcast(podcast *models.Podcast) error {
	filePath := fmt.Sprintf("%s/podcasts/%s.json", fs.dataDir, podcast.ID)
	// Ensure the directory exists

	// Marshal the podcast data to JSON

	data, err := json.MarshalIndent(podcast, "", "  ")
	if err != nil {
		return err
	}
	os.Create(filePath) // Ensure the directory exists
	return os.WriteFile(filePath, data, 0644)
}

func (fs *FileStorage) LoadPodcast(id string) (*models.Podcast, error) {
	filePath := fmt.Sprintf("%s/podcasts/%s.json", fs.dataDir, id)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var podcast models.Podcast
	err = json.Unmarshal(data, &podcast)
	if err != nil {
		return nil, err
	}

	return &podcast, nil
}

func (fs *FileStorage) SaveEpisode(episode *models.Episode) error {
	episodeId := strings.Replace(episode.ID, "/", "_", -1) // Replace spaces with underscores for file naming
	filePath := fmt.Sprintf("%s/episodes/%s/%s.json", fs.dataDir, episode.PodcastID, episodeId)
	// Ensure the directory exists
	if err := os.MkdirAll(fmt.Sprintf("%s/episodes/%s", fs.dataDir, episode.PodcastID), 0755); err != nil {
		return fmt.Errorf("failed to create directory for episodes: %w", err)
	}
	data, err := json.MarshalIndent(episode, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
func (fs *FileStorage) LoadEpisode(podcastID, episodeID string) (*models.Episode, error) {
	episodeID = strings.Replace(episodeID, "/", "_", -1) // Replace spaces with underscores for file naming
	filePath := fmt.Sprintf("%s/episodes/%s/%s.json", fs.dataDir, podcastID, episodeID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var episode models.Episode
	err = json.Unmarshal(data, &episode)
	if err != nil {
		return nil, err
	}

	return &episode, nil
}
