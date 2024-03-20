package backgroundjob

import (
	"context"
	"log"
	"sync"
	"time"
)

type BackgroundJob struct{}

func New() *BackgroundJob {
	return &BackgroundJob{}
}

type JobOptions struct {
	Interval time.Duration

	// DelayedStart has more priority than StartAfter and use Interval
	DelayedStart *bool
	StartAfter   *time.Duration

	Job func()
}

// Run is eminently start goroutine with callback
func (b *BackgroundJob) Run(options JobOptions, ctx context.Context) {
	go func() {
		if options.DelayedStart != nil {
			if *options.DelayedStart {
				time.Sleep(options.Interval)
			}
		} else {
			if options.StartAfter != nil {
				time.Sleep(*options.StartAfter)
			}
		}

		b.job(options.Interval, options.Job, ctx)
	}()
}

func (b *BackgroundJob) job(interval time.Duration, job func(), ctx context.Context) {
	for {
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			job()
			defer wg.Done()
		}()

		wg.Wait()

		select {
		case <-ctx.Done():
			log.Println("background job done")
			return
		default:
		}

		time.Sleep(interval)
	}
}
