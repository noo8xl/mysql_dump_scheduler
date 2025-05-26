package common

import (
	"os"
	"time"
)

type DatabaseConfig struct {
	Host         string    `json:"host"`
	Port         string    `json:"port"`
	User         string    `json:"user"`
	Password     string    `json:"password"`
	Database     string    `json:"database"`
	SqlFilesPath *SqlFiles `json:"file_path"` // path to the file to insert the data
}

type TelegramConfig struct {
	ChatId string `json:"chat_id"` // telegram user chat id
	Token  string `json:"token"`   // telegram bot token
}

type SqlFiles struct {
	TablesFilePath string `json:"tables_file_path"` // path to the tables file
	DataFilePath   string `json:"data_file_path"`   // path to the data file
	DumpDirPath    string `json:"dump_dir_path"`    // expected path to the dump file
}

type MakeOpts struct {
	RunPath string
}

type SchedulerConfig struct {
	File     *os.File      `json:"file"`      // file will be set after the dump is created
	Duration time.Duration `json:"duration"`  // duration of the scheduler
	MakeOpts MakeOpts      `json:"make_opts"` // makefile options
}
