package taskrunner

import (
	"context"
	"fmt"
	"monitor/internal/config"
	"monitor/internal/services/check"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func Launch(ctx context.Context, c config.Check) error {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			time.Duration(c.Interval)*time.Second,
		),
		gocron.NewTask(
			func(ctx context.Context, c config.Check) {
				check.Run(ctx, c)
			},
			ctx,
			c,
		),
	)
	if err != nil {
		return err
	}
	// each job has a unique id
	fmt.Printf("Startup: %v\n", j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-ctx.Done():
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		return err
	}
	fmt.Printf("Shutdown: %v\n", j.ID())

	return nil
}
