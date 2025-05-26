package scheduler

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/noo8xl/mysql_dump_scheduler/common"
)

type SchedulerService struct {
	DatabaseConfig  *common.DatabaseConfig
	TelegramConfig  *common.TelegramConfig
	SchedulerConfig *common.SchedulerConfig
	logChan         chan string
}

// @InitScheduler -> Initialize the scheduler
func InitSchedulerService() *SchedulerService {
	return &SchedulerService{
		logChan: make(chan string, 500),
	}
}

func (s *SchedulerService) SetDatabaseConfig(opts *common.DatabaseConfig) error {

	if opts.Host == "" || opts.Port == "" || opts.User == "" || opts.Password == "" || opts.Database == "" {
		return errors.New("error: databaseConfig is not set")
	}

	s.DatabaseConfig = opts
	return nil
}

func (s *SchedulerService) SetTelegramConfig(opts *common.TelegramConfig) error {

	if opts.ChatId == "" {
		return errors.New("error: telegramConfig is not set")
	}

	s.TelegramConfig = opts
	return nil
}

func (s *SchedulerService) SetSchedulerConfig(opts *common.SchedulerConfig) error {

	if opts.Duration == 0 {
		return errors.New("error: duration is not set")
	}

	s.SchedulerConfig = opts
	return nil
}

// @Bootstrap -> bootstrap the scheduler
func (s *SchedulerService) Bootstrap(ctx context.Context) error {

	ticker := time.NewTicker(s.SchedulerConfig.Duration)
	defer ticker.Stop()

	s.logChan <- "SchedulerService::Bootstrap: started"

	for {
		select {
		case <-ticker.C:
			if err := s.run(); err != nil {
				return err
			}
		case msg := <-s.logChan:
			log.Println(msg)
		case <-ctx.Done():
			s.logChan <- "Scheduler shutting down"
			return nil
		}
	}
}
