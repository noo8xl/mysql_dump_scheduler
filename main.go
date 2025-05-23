package scheduler

import (
	"context"
	"errors"
	"os"
	"time"
)

type SchedulerConfig struct {
	Path     string        `json:"path"`     // path to the file
	File     *os.File      `json:"file"`     // file will be set after the dump is created
	Duration time.Duration `json:"duration"` // duration of the scheduler
	// Timer    *time.Timer   `json:"timer"`    // timer will be set after the duration is set
}

type TelegramConfig struct {
	ChatId string `json:"chat_id"` // telegram user chat id
	Token  string `json:"token"`   // telegram bot token
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type SchedulerOpts struct {
	DatabaseConfig  DatabaseConfig
	TelegramConfig  TelegramConfig
	SchedulerConfig SchedulerConfig
}

type SchedulerService struct {
	DatabaseConfig  DatabaseConfig
	TelegramConfig  TelegramConfig
	SchedulerConfig SchedulerConfig
}

// @InitScheduler -> Initialize the scheduler
func InitScheduler(opts SchedulerOpts) (*SchedulerService, error) {

	if opts.DatabaseConfig.Host == "" || opts.DatabaseConfig.Port == "" || opts.DatabaseConfig.User == "" || opts.DatabaseConfig.Password == "" || opts.DatabaseConfig.Name == "" {
		return nil, errors.New("error: DatabaseConfig is not set")
	}

	if opts.TelegramConfig.ChatId == "" {
		return nil, errors.New("error: TelegramConfig is not set")
	}

	if opts.SchedulerConfig.Duration == 0 {
		return nil, errors.New("error: Duration is not set")
	}

	return &SchedulerService{
		DatabaseConfig:  opts.DatabaseConfig,
		TelegramConfig:  opts.TelegramConfig,
		SchedulerConfig: opts.SchedulerConfig,
	}, nil
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

// @SendFile -> send the file to the telegram
func (s *SchedulerService) SendFile() error {

	compressedFile, err := s.getCompressedFile()
	if err != nil {
		return err
	}

	if err := s.sendFileToTelegram(compressedFile); err != nil {
		return err
	}

	return nil
}
