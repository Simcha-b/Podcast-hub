package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Simcha-b/Podcast-Hub/models"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	SavePodcast(podcast *models.Podcast) error
	LoadAllPodcasts() ([]models.Podcast, error)
	LoadPodcastByID(id string) (*models.Podcast, error)
	SaveEpisode(episode *models.Episode) error
	LoadEpisodes(podcastID string) ([]models.Episode, error)
	LoadEpisodeByID(podcastID, episodeID string) (*models.Episode, error)
	SearchPodcasts(query string) ([]models.Podcast, error)
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

	
	//check if the podcast exist 
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
	}
	Logger.Error(fmt.Sprintf("Podcast with ID %s not found", id))
	return nil, ErrNotFound
}

func (fs *FileStorage) SaveEpisode(episode *models.Episode) error {
	episodesPath := fmt.Sprintf("%s/episodes/episodes_%s.json", fs.dataDir, episode.PodcastID)
	episodes, err := fs.LoadEpisodes(episode.PodcastID)
	if err != nil && err != ErrNotFound {
		return fmt.Errorf("failed to load episodes for podcast %s: %w", episode.PodcastID, err)
	}

	// אם הקובץ לא קיים (ErrNotFound), תתחיל עם רשימה ריקה
	if err == ErrNotFound {
		episodes = []models.Episode{}
	}

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
		return nil, ErrNotFound
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
			Logger.Info(fmt.Sprintf("Episode with ID %s found in podcast %s", episodeID, podcastID))
			return &ep, nil
		}
	}
	Logger.Info(fmt.Sprintf("Episode with ID %s not found in podcast %s", episodeID, podcastID))
	return nil, ErrNotFound
}

func (fs *FileStorage) LoadLastEpisodes() ([]models.Episode, error) {
	// Check if the episodes directory exists
	files, err := os.ReadDir(fmt.Sprintf("%s/episodes", fs.dataDir))
	if err != nil {
		return nil, fmt.Errorf("failed to read episodes directory: %w", err)
	}

	// Filter relevant files first
	var episodeFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "episodes_") && strings.HasSuffix(file.Name(), ".json") {
			episodeFiles = append(episodeFiles, file)
		}
	}

	if len(episodeFiles) == 0 {
		Logger.Info("No episode files found")
		return nil, fmt.Errorf("no episode files found")
	}

	// Channel for collecting results
	type result struct {
		episodes []models.Episode
		err      error
		filename string
	}

	resultChan := make(chan result, len(episodeFiles))

	// Process files concurrently
	var wg sync.WaitGroup
	weekAgo := time.Now().AddDate(0, 0, -7)

	for _, file := range episodeFiles {
		wg.Add(1)
		go func(f os.DirEntry) {
			defer wg.Done()

			filePath := filepath.Join(fs.dataDir, "episodes", f.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				resultChan <- result{
					err:      fmt.Errorf("failed to read episode file %s: %w", f.Name(), err),
					filename: f.Name(),
				}
				return
			}

			Logger.Info(fmt.Sprintf("Processing file: %s", f.Name()))

			var eps []models.Episode
			err = json.Unmarshal(data, &eps)
			if err != nil {
				resultChan <- result{
					err:      fmt.Errorf("failed to unmarshal episodes from file %s: %w", f.Name(), err),
					filename: f.Name(),
				}
				return
			}

			// Filter episodes from last 7 days
			var recentEpisodes []models.Episode
			for i := range eps {
				if eps[i].PublishedAt.After(weekAgo) {
					recentEpisodes = append(recentEpisodes, eps[i])
				}
			}

			resultChan <- result{
				episodes: recentEpisodes,
				filename: f.Name(),
			}
		}(file)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var lasteps []models.Episode
	for res := range resultChan {
		if res.err != nil {
			return nil, res.err
		}
		lasteps = append(lasteps, res.episodes...)
	}

	if len(lasteps) == 0 {
		Logger.Info("No episodes found in the last 7 days")
		return nil, fmt.Errorf("no episodes found in the last 7 days")
	}

	Logger.Info(fmt.Sprintf("Found %d episodes in the last 7 days", len(lasteps)))
	return lasteps, nil
}

func (fs *FileStorage) SearchPodcasts(query string) ([]models.Podcast, error) {
	podcasts, err := fs.LoadAllPodcasts()
	if err != nil {
		return nil, fmt.Errorf("failed to load podcasts: %w", err)
	}

	var results []models.Podcast
	for _, podcast := range podcasts {
		if containsIgnoreCase(podcast.Title, query) || containsIgnoreCase(podcast.Description, query) {
			results = append(results, podcast)
		}
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}
	return results, nil
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
