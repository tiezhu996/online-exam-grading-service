package worker

import (
	"context"
	"sync"

	"onlineexam/internal/model"
	"onlineexam/internal/service"
)

type Pool struct {
	svc     *service.Service
	workers int
}

func New(svc *service.Service, workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}
	return &Pool{svc: svc, workers: workers}
}

func (p *Pool) GradeAll(ctx context.Context) model.Summary {
	batches := p.svc.SubmissionBatches()

	var wg sync.WaitGroup
	ch := make(chan []*model.Submission, len(batches))

	go func() {
		defer close(ch)
		for _, b := range batches {
			select {
			case <-ctx.Done():
				return
			case ch <- b:
			}
		}
	}()

	var mu sync.Mutex
	var sum model.Summary

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range ch {
				var local model.Summary
				for _, sub := range batch {
					select {
					case <-ctx.Done():
						return
					default:
					}
					_, err := p.svc.GradeSubmission(sub.ID)
					local.Checked++
					if err != nil {
						local.Failed++
						continue
					}
					local.Graded++
				}
				mu.Lock()
				sum = model.MergeSummary(sum, local)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return sum
}
