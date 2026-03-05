package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// GenerationRequest represents a queued document generation request.
type GenerationRequest struct {
	GenerationID string
	TemplateID   string
	TenantID     string
	EntityID     string
	EntityType   string
	Format       string                 // pdf, html
	Data         map[string]interface{}
	ActorID      string
}

// GenerationWorker manages asynchronous document generation with a worker pool.
type GenerationWorker struct {
	queue     chan GenerationRequest
	workers   int
	generator func(ctx context.Context, req GenerationRequest) error
	wg        sync.WaitGroup
}

// NewGenerationWorker creates a new document generation worker pool.
// workerCount specifies the number of concurrent workers.
// queueSize specifies the capacity of the work queue.
// generator is the function that performs the actual document generation.
func NewGenerationWorker(
	workerCount int,
	queueSize int,
	generator func(ctx context.Context, req GenerationRequest) error,
) *GenerationWorker {
	if workerCount <= 0 {
		workerCount = 1
	}
	if queueSize <= 0 {
		queueSize = 100
	}
	return &GenerationWorker{
		queue:     make(chan GenerationRequest, queueSize),
		workers:   workerCount,
		generator: generator,
	}
}

// Start initiates the worker goroutines.
// Each worker pulls requests from the queue and processes them.
func (w *GenerationWorker) Start(ctx context.Context) {
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go w.processWorker(ctx, i)
	}
}

// Enqueue adds a generation request to the queue.
// Returns an error if the queue is full (non-blocking).
func (w *GenerationWorker) Enqueue(req GenerationRequest) error {
	select {
	case w.queue <- req:
		return nil
	default:
		return fmt.Errorf("generation queue is full")
	}
}

// Stop signals all workers to stop and waits for them to finish.
func (w *GenerationWorker) Stop() {
	close(w.queue)
	w.wg.Wait()
}

// processWorker is the main worker loop.
func (w *GenerationWorker) processWorker(ctx context.Context, workerID int) {
	defer w.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("worker %d panicked: %v", workerID, r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("[docgen-worker-%d] shutting down due to context cancellation", workerID)
			return
		case req, ok := <-w.queue:
			if !ok {
				logger.Infof("[docgen-worker-%d] shutting down, queue closed", workerID)
				return
			}
			w.processRequest(ctx, req, workerID)
		}
	}
}

// processRequest handles a single generation request.
func (w *GenerationWorker) processRequest(ctx context.Context, req GenerationRequest, workerID int) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[docgen-worker-%d] panic processing generation %s: %v",
				workerID, req.GenerationID, r)
		}
	}()

	if w.generator == nil {
		logger.Errorf("[docgen-worker-%d] generator function is nil for generation %s",
			workerID, req.GenerationID)
		return
	}

	logger.Infof("[docgen-worker-%d] starting generation id=%s template=%s tenant=%s format=%s",
		workerID, req.GenerationID, req.TemplateID, req.TenantID, req.Format)

	if err := w.generator(ctx, req); err != nil {
		logger.Errorf("[docgen-worker-%d] generation failed id=%s: %v",
			workerID, req.GenerationID, err)
		return
	}

	logger.Infof("[docgen-worker-%d] generation completed id=%s", workerID, req.GenerationID)
}
