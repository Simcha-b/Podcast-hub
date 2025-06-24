 package services

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/Simcha-b/Podcast-Hub/config"
// 	"github.com/Simcha-b/Podcast-Hub/models"
// 	"golang.org/x/time/rate"
// )

// // FeedJob represents a job to process a single feed
// type FeedJob struct {
// 	Feed    models.Feed
// 	Storage *FileStorage
// }

// // FeedResult represents the result of processing a feed
// type FeedResult struct {
// 	Feed    models.Feed
// 	Success bool
// 	Error   error
// 	Stats   ProcessingStats
// }

// // ProcessingStats holds statistics for feed processing
// type ProcessingStats struct {
// 	PodcastsProcessed int64
// 	EpisodesProcessed int64
// 	NewEpisodes       int64
// 	ProcessingTime    time.Duration
// }

// // WorkerPool manages concurrent feed processing with rate limiting
// type WorkerPool struct {
// 	config      *config.Config
// 	jobs        chan FeedJob
// 	results     chan FeedResult
// 	workers     []*Worker
// 	rateLimiter *rate.Limiter
// 	stats       *GlobalStats
// 	ctx         context.Context
// 	cancel      context.CancelFunc
// 	wg          sync.WaitGroup
// }

// // GlobalStats holds global processing statistics
// type GlobalStats struct {
// 	TotalJobs       int64
// 	CompletedJobs   int64
// 	FailedJobs      int64
// 	TotalPodcasts   int64
// 	TotalEpisodes   int64
// 	TotalNewEpisodes int64
// 	mu              sync.RWMutex
// }

// // Worker represents a single worker in the pool
// type Worker struct {
// 	id          int
// 	pool        *WorkerPool
// 	rateLimiter *rate.Limiter
// }

// // NewWorkerPool creates a new worker pool with the given configuration
// func NewWorkerPool(cfg *config.Config) *WorkerPool {
// 	ctx, cancel := context.WithCancel(context.Background())
	
// 	// Create rate limiter: cfg.RATE_LIMIT requests per second
// 	rateLimiter := rate.NewLimiter(rate.Limit(cfg.RATE_LIMIT), cfg.RATE_LIMIT)
	
// 	pool := &WorkerPool{
// 		config:      cfg,
// 		jobs:        make(chan FeedJob, cfg.CHANNEL_BUFFER_SIZE),
// 		results:     make(chan FeedResult, cfg.CHANNEL_BUFFER_SIZE),
// 		rateLimiter: rateLimiter,
// 		stats:       &GlobalStats{},
// 		ctx:         ctx,
// 		cancel:      cancel,
// 	}
	
// 	// Create workers
// 	pool.workers = make([]*Worker, cfg.MAX_WORKERS)
// 	for i := 0; i < cfg.MAX_WORKERS; i++ {
// 		pool.workers[i] = &Worker{
// 			id:          i,
// 			pool:        pool,
// 			rateLimiter: rate.NewLimiter(rate.Limit(cfg.RATE_LIMIT), 1),
// 		}
// 	}
	
// 	return pool
// }

// // Start starts all workers in the pool
// func (wp *WorkerPool) Start() {
// 	Logger.Info(fmt.Sprintf("Starting worker pool with %d workers", wp.config.MAX_WORKERS))
	
// 	for _, worker := range wp.workers {
// 		wp.wg.Add(1)
// 		go worker.start()
// 	}
	
// 	// Start result processor
// 	wp.wg.Add(1)
// 	go wp.processResults()
// }

// // Stop gracefully stops the worker pool
// func (wp *WorkerPool) Stop() {
// 	Logger.Info("Stopping worker pool...")
	
// 	// Close jobs channel to signal no more jobs
// 	close(wp.jobs)
	
// 	// Cancel context to stop workers
// 	wp.cancel()
	
// 	// Wait for all workers to finish
// 	wp.wg.Wait()
	
// 	// Close results channel
// 	close(wp.results)
	
// 	Logger.Info("Worker pool stopped")
// }

// // SubmitJob submits a job to the worker pool
// func (wp *WorkerPool) SubmitJob(job FeedJob) {
// 	select {
// 	case wp.jobs <- job:
// 		atomic.AddInt64(&wp.stats.TotalJobs, 1)
// 	case <-wp.ctx.Done():
// 		Logger.Error("Cannot submit job: worker pool is shutting down")
// 	}
// }

// // GetStats returns a copy of the current global statistics
// func (wp *WorkerPool) GetStats() GlobalStats {
// 	wp.stats.mu.RLock()
// 	defer wp.stats.mu.RUnlock()
	
// 	return GlobalStats{
// 		TotalJobs:        atomic.LoadInt64(&wp.stats.TotalJobs),
// 		CompletedJobs:    atomic.LoadInt64(&wp.stats.CompletedJobs),
// 		FailedJobs:       atomic.LoadInt64(&wp.stats.FailedJobs),
// 		TotalPodcasts:    atomic.LoadInt64(&wp.stats.TotalPodcasts),
// 		TotalEpisodes:    atomic.LoadInt64(&wp.stats.TotalEpisodes),
// 		TotalNewEpisodes: atomic.LoadInt64(&wp.stats.TotalNewEpisodes),
// 	}
// }

// // start starts a single worker
// func (w *Worker) start() {
// 	defer w.pool.wg.Done()
	
// 	Logger.Info(fmt.Sprintf("Worker %d started", w.id))
	
// 	for {
// 		select {
// 		case job, ok := <-w.pool.jobs:
// 			if !ok {
// 				Logger.Info(fmt.Sprintf("Worker %d: jobs channel closed, stopping", w.id))
// 				return
// 			}
			
// 			// Process the job with rate limiting and timeout
// 			result := w.processJob(job)
			
// 			// Send result
// 			select {
// 			case w.pool.results <- result:
// 			case <-w.pool.ctx.Done():
// 				return
// 			}
			
// 		case <-w.pool.ctx.Done():
// 			Logger.Info(fmt.Sprintf("Worker %d: context cancelled, stopping", w.id))
// 			return
// 		}
// 	}
// }

// // processJob processes a single feed job with rate limiting and timeout
// func (w *Worker) processJob(job FeedJob) FeedResult {
// 	startTime := time.Now()
	
// 	// Apply rate limiting
// 	if err := w.rateLimiter.Wait(w.pool.ctx); err != nil {
// 		return FeedResult{
// 			Feed:    job.Feed,
// 			Success: false,
// 			Error:   fmt.Errorf("rate limiting error: %w", err),
// 		}
// 	}
	
// 	// Create context with timeout for this specific job
// 	jobCtx, cancel := context.WithTimeout(w.pool.ctx, w.pool.config.PROCESSING_TIMEOUT)
// 	defer cancel()
	
// 	// Process the feed with timeout
// 	result := make(chan FeedResult, 1)
// 	go func() {
// 		podcast, episodes, err := parseRSSFeedWithTimeout(job.Feed.URL, w.pool.config.REQUEST_TIMEOUT)
// 		if err != nil {
// 			result <- FeedResult{
// 				Feed:    job.Feed,
// 				Success: false,
// 				Error:   err,
// 				Stats: ProcessingStats{
// 					ProcessingTime: time.Since(startTime),
// 				},
// 			}
// 			return
// 		}
		
// 		// Save podcast and episodes
// 		stats := ProcessingStats{
// 			ProcessingTime: time.Since(startTime),
// 		}
		
// 		if podcast != nil {
// 			if err := job.Storage.SavePodcast(podcast); err != nil {
// 				result <- FeedResult{
// 					Feed:    job.Feed,
// 					Success: false,
// 					Error:   fmt.Errorf("failed to save podcast: %w", err),
// 					Stats:   stats,
// 				}
// 				return
// 			}
// 			stats.PodcastsProcessed = 1
// 		}
		
// 		// Process episodes
// 		newCount, err := w.saveEpisodesWithDuplicateCheck(job.Storage, podcast.ID, episodes)
// 		if err != nil {
// 			result <- FeedResult{
// 				Feed:    job.Feed,
// 				Success: false,
// 				Error:   fmt.Errorf("failed to save episodes: %w", err),
// 				Stats:   stats,
// 			}
// 			return
// 		}
		
// 		stats.EpisodesProcessed = int64(len(episodes))
// 		stats.NewEpisodes = int64(newCount)
		
// 		result <- FeedResult{
// 			Feed:    job.Feed,
// 			Success: true,
// 			Error:   nil,
// 			Stats:   stats,
// 		}
// 	}()
	
// 	// Wait for result or timeout
// 	select {
// 	case res := <-result:
// 		return res
// 	case <-jobCtx.Done():
// 		return FeedResult{
// 			Feed:    job.Feed,
// 			Success: false,
// 			Error:   fmt.Errorf("job timeout after %v", w.pool.config.PROCESSING_TIMEOUT),
// 			Stats: ProcessingStats{
// 				ProcessingTime: time.Since(startTime),
// 			},
// 		}
// 	}
// }

// // saveEpisodesWithDuplicateCheck saves episodes while checking for duplicates
// func (w *Worker) saveEpisodesWithDuplicateCheck(storage *FileStorage, podcastID string, episodes []models.Episode) (int, error) {
// 	existingEpisodes, err := storage.LoadEpisodes(podcastID)
// 	if err != nil && !isFileNotExistError(err) {
// 		return 0, err
// 	}
	
// 	existingMap := make(map[string]models.Episode)
// 	for _, ep := range existingEpisodes {
// 		existingMap[ep.ID] = ep
// 	}
	
// 	newCount := 0
// 	for _, episode := range episodes {
// 		existing, ok := existingMap[episode.ID]
// 		if !ok || episode.PublishedAt.After(existing.PublishedAt) {
// 			if err := storage.SaveEpisode(&episode); err != nil {
// 				return newCount, err
// 			}
// 			newCount++
// 		}
// 	}
	
// 	return newCount, nil
// }

// // processResults processes results from workers and updates global statistics
// func (wp *WorkerPool) processResults() {
// 	defer wp.wg.Done()
	
// 	for {
// 		select {
// 		case result, ok := <-wp.results:
// 			if !ok {
// 				Logger.Info("Results channel closed, stopping result processor")
// 				return
// 			}
			
// 			// Update global statistics
// 			if result.Success {
// 				atomic.AddInt64(&wp.stats.CompletedJobs, 1)
// 				atomic.AddInt64(&wp.stats.TotalPodcasts, result.Stats.PodcastsProcessed)
// 				atomic.AddInt64(&wp.stats.TotalEpisodes, result.Stats.EpisodesProcessed)
// 				atomic.AddInt64(&wp.stats.TotalNewEpisodes, result.Stats.NewEpisodes)
				
// 				Logger.Info(fmt.Sprintf("Successfully processed feed %s: %d new episodes in %v", 
// 					result.Feed.URL, result.Stats.NewEpisodes, result.Stats.ProcessingTime))
				
// 				// Update feed status
// 				result.Feed.LastFetched = time.Now()
// 				if err := UpdateFeedStatus(result.Feed, true); err != nil {
// 					Logger.Error(fmt.Sprintf("Failed to update feed status for %s: %v", result.Feed.URL, err))
// 				}
// 			} else {
// 				atomic.AddInt64(&wp.stats.FailedJobs, 1)
// 				Logger.Error(fmt.Sprintf("Failed to process feed %s: %v", result.Feed.URL, result.Error))
				
// 				// Update feed status
// 				if err := UpdateFeedStatus(result.Feed, false); err != nil {
// 					Logger.Error(fmt.Sprintf("Failed to update feed status for %s: %v", result.Feed.URL, err))
// 				}
// 			}
			
// 		case <-wp.ctx.Done():
// 			Logger.Info("Context cancelled, stopping result processor")
// 			return
// 		}
// 	}
// }