package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/noo8xl/mysql_dump_scheduler/common"
)

type SchedulerService struct {
	DatabaseConfig  *common.DatabaseConfig
	TelegramConfig  *common.TelegramConfig
	SchedulerConfig *common.SchedulerConfig
}

// @InitScheduler -> Initialize the scheduler
func InitScheduler() *SchedulerService {
	return &SchedulerService{}
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

	for {
		select {
		case <-ticker.C:
			if err := s.Run(); err != nil {
				return err
			}
			if s.TelegramConfig.ChatId != "" {
				if err := s.sendFileToTelegram(); err != nil {
					return err
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// @Run -> run the scheduler in following order:
//  1. set the timer
//  2. create the dump file
//  3. compress the dump file
//  4. send the file to the telegram // -> optional
func (s *SchedulerService) Run() error {

	file, err := s.createDumpFile()
	if err != nil {
		return err
	}

	s.SchedulerConfig.File = file
	if err := s.compressDumpFile(); err != nil {
		return err
	}

	return nil
}
