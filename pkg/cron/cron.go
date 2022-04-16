package cron

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type CronService struct{}

func (cs *CronService) Serve(_ context.Context) error {
	s := gocron.NewScheduler(time.UTC)
	// s.Every("30m").Do(func() {
	// 	log.Info("Starting scheduled task for ingestion")
	// 	if count, err := ingestor.IngestFromAll(); err != nil {
	// 		log.Error("ingestor failed: %v", err)
	// 	} else {
	// 		log.Infof("ingestor completed. processed %d rows.", count)
	// 	}
	// 	log.Infof("Finished scheduled task")
	// })

	log.Info("Starting go-cron scheduler.")
	s.StartBlocking()
	return nil
}
