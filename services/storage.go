package services

import (
	"encoding/json"
	"fmt"
	"os"
	. "github.com/Simcha-b/Podcast-Hub/models"
)

type Storage interface {
	SavePodcast(podcast *Podcast) error
	LoadPodcast(id string) (*Podcast, error)
	SaveEpisode(episode *Episode) error
	LoadEpisode(podcastID, episodeID string) (*Episode, error)
}

type FileStorage struct {
	dataDir string
}

func NewFileStorage(dataDir string) *FileStorage {
	return &FileStorage{
		dataDir: dataDir,
	}
}

func (fs *FileStorage) SavePodcast(podcast *Podcast) error {
	filePath := fmt.Sprintf("../%s/podcasts/%s.json", fs.dataDir, podcast.ID)
	// Ensure the directory exists

	data, err := json.MarshalIndent(podcast, "", "  ")
	if err != nil {
		return err
	}
	os.Create(filePath) // Ensure the directory exists
	return os.WriteFile(filePath, data, 0644)
}

func (fs *FileStorage) LoadPodcast(id string) (*Podcast, error) {
	filePath := fmt.Sprintf("../%s/podcasts/%s.json", fs.dataDir, id)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var podcast Podcast
	err = json.Unmarshal(data, &podcast)
	if err != nil {
		return nil, err
	}

	return &podcast, nil
}

func (fs *FileStorage) SaveEpisode(episode *Episode) error {
	filePath := fmt.Sprintf("../%s/episodes/%s/%s.json", fs.dataDir, episode.PodcastID, episode.ID)
	data, err := json.MarshalIndent(episode, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
func (fs *FileStorage) LoadEpisode(podcastID, episodeID string) (*Episode, error) {
	filePath := fmt.Sprintf("../%s/episodes/%s/%s.json", fs.dataDir, podcastID, episodeID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var episode Episode
	err = json.Unmarshal(data, &episode)
	if err != nil {
		return nil, err
	}

	return &episode, nil
}

// func main() {
// 	// Example usage
// 	storage := NewFileStorage("data")

// 	podcast := &Podcast{
// 		ID:          "podcast1",
// 		Title:       "My Podcast",
// 		Description: "A great podcast about interesting topics.",
// 		ImageURL:    "http://example.com/image.jpg",
// 	}

// 	err := storage.SavePodcast(podcast)
// 	if err != nil {
// 		fmt.Println("Error saving podcast:", err)
// 		return
// 	}

// 	loadedPodcast, err := storage.LoadPodcast("podcast1")
// 	if err != nil {
// 		fmt.Println("Error loading podcast:", err)
// 		return
// 	}

// 	fmt.Printf("Loaded podcast: %+v\n", loadedPodcast)
// }
