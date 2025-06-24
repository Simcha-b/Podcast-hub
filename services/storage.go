package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Simcha-b/Podcast-Hub/models"
)

type Storage interface {
	SavePodcast(podcast *models.Podcast) error
	LoadAllPodcasts() ([]models.Podcast, error)
	LoadPodcastByID(id string) (*models.Podcast, error)
	SaveEpisode(episode *models.Episode) error
	LoadEpisodes(podcastID string) ([]models.Episode, error)
	LoadEpisodeByID(podcastID, episodeID string) (*models.Episode, error)
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
	allPodcastsPath := fmt.Sprintf("%s/podcasts/all_podcasts.json", fs.dataDir)
	podcasts, err := fs.LoadAllPodcasts()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load all podcasts: %w", err)
	}

	// בדיקה אם הפודקאסט כבר קיים (עדכון במקום הוספה כפולה)
	updated := false
	for i, p := range podcasts {
		if p.ID == podcast.ID {
			podcasts[i] = *podcast
			updated = true
			break
		}
	}
	if !updated {
		podcasts = append(podcasts, *podcast)
	}

	data, err := json.MarshalIndent(podcasts, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/podcasts", fs.dataDir), 0755); err != nil {
		return fmt.Errorf("failed to create directory for podcasts: %w", err)
	}

	return os.WriteFile(allPodcastsPath, data, 0644)
}

func (fs *FileStorage) LoadAllPodcasts() ([]models.Podcast, error) {
	allPodcastsPath := fmt.Sprintf("%s/podcasts/all_podcasts.json", fs.dataDir)
	data, err := os.ReadFile(allPodcastsPath)
	if err != nil {
		return nil, err
	}

	var podcasts []models.Podcast
	err = json.Unmarshal(data, &podcasts)
	if err != nil {
		return nil, err
	}

	return podcasts, nil
}

func (fs *FileStorage) LoadPodcastByID(id string) (*models.Podcast, error) {
	podcasts, err := fs.LoadAllPodcasts()
	if err != nil {
		return nil, fmt.Errorf("failed to load all podcasts: %w", err)
	}
	for _, podcast := range podcasts {
		if podcast.ID == id {
			Logger.Info(fmt.Sprintf("Podcast with ID %s found", id))
			return &podcast, nil
		}
		Logger.Info(fmt.Sprintf("Podcast with ID %s not found", id))
	}
	return nil, nil
}

func (fs *FileStorage) SaveEpisode(episode *models.Episode) error {
	episodesPath := fmt.Sprintf("%s/episodes/episodes_%s.json", fs.dataDir, episode.PodcastID)
	episodes, err := fs.LoadEpisodes(episode.PodcastID)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load episodes for podcast %s: %w", episode.PodcastID, err)
	}

	// בדיקה אם הפרק כבר קיים (עדכון במקום הוספה כפולה)
	updated := false
	for i, ep := range episodes {
		if ep.ID == episode.ID {
			episodes[i] = *episode
			updated = true
			break
		}
	}
	if !updated {
		episodes = append(episodes, *episode)
	}

	data, err := json.MarshalIndent(episodes, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/episodes", fs.dataDir), 0755); err != nil {
		return fmt.Errorf("failed to create directory for episodes: %w", err)
	}

	return os.WriteFile(episodesPath, data, 0644)
}

func (fs *FileStorage) LoadEpisodes(podcastID string) ([]models.Episode, error) {
	episodesPath := fmt.Sprintf("%s/episodes/episodes_%s.json", fs.dataDir, podcastID)
	data, err := os.ReadFile(episodesPath)
	if err != nil {
		return nil, err
	}

	var episodes []models.Episode
	err = json.Unmarshal(data, &episodes)
	if err != nil {
		return nil, err
	}

	return episodes, nil
}

func (fs *FileStorage) LoadEpisodeByID(podcastID, episodeID string) (*models.Episode, error) {
	episodes, err := fs.LoadEpisodes(podcastID)
	if err != nil {
		return nil, err
	}

	for _, ep := range episodes {
		if ep.ID == episodeID {
			return &ep, nil
		}
	}
	return nil, fmt.Errorf("episode with ID %s for podcast %s not found", episodeID, podcastID)
}
