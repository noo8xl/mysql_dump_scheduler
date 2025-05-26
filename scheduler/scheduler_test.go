package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/noo8xl/mysql_dump_scheduler/common"
	"github.com/noo8xl/mysql_dump_scheduler/scheduler"
)

var (
	schedSvc *scheduler.SchedulerService

	dbConfig = &common.DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "password",
		Database: "test",
		SqlFilesPath: &common.SqlFiles{
			TablesFilePath: "path/to/tables/db.sql",
			DataFilePath:   "path/to/data/data.sql",
			DumpDirPath:    "path/to/dump/dir",
		},
	}

	telegramConfig = &common.TelegramConfig{
		ChatId: "ur-chat-id-here",
		Token:  "ur-token-here",
	}

	schedulerConfig = &common.SchedulerConfig{
		Duration: 3 * time.Second,
		MakeOpts: &common.MakeOpts{
			RunPath: "path/to/dir/with/your/Makefile",
		},
	}
)

func init() {
	schedSvc = scheduler.InitSchedulerService()
}

func TestSetDatabaseConfig(t *testing.T) {
	if err := schedSvc.SetDatabaseConfig(dbConfig); err != nil {
		t.Fatalf("error: failed to set database config: %v", err)
	}
}

func TestSetTelegramConfig(t *testing.T) {
	if err := schedSvc.SetTelegramConfig(telegramConfig); err != nil {
		t.Fatalf("error: failed to set telegram config: %v", err)
	}
}

func TestSetSchedulerConfig(t *testing.T) {
	if err := schedSvc.SetSchedulerConfig(schedulerConfig); err != nil {
		t.Fatalf("error: failed to set scheduler config: %v", err)
	}
}

func TestBootstrap(t *testing.T) {
	if err := schedSvc.Bootstrap(context.Background()); err != nil {
		t.Fatalf("error: failed to bootstrap: %v", err)
	}
}
