package main

import (
	"context"
	"log"
	"monitor/internal/config"
	"monitor/internal/services/taskrunner"
	"monitor/internal/services/workcalendar"

	"monitor/internal/utils"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const shutdownTimeoutSecs = 10

func main() {
	cts := InitShutdownHandler(context.Background(), shutdownTimeoutSecs*time.Second)

	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalln(err)
	}

	cfg.CalendarMap = workcalendar.InitCalendarService(cfg.Calendars)
	if err != nil {
		log.Fatalln(err)
	}
	//type ky string

	ctx := context.WithValue(cts, "config", &cfg)

	utils.Dump("Config", cfg)
	if err != nil {
		log.Fatalln(err)
	}
	g, ctx := errgroup.WithContext(ctx)
	for _, ch := range cfg.Checks {
		g.Go(func() error {
			err := taskrunner.Launch(ctx, ch)
			return err
		})
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func InitShutdownHandler(ctx context.Context, shutdownTimeout time.Duration) context.Context {
	newCtx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-newCtx.Done()
		stop() // don't handle subsequent signals

		timeoutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		<-timeoutCtx.Done()

		if err := timeoutCtx.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	return newCtx
}
