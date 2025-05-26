package initializers

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/noo8xl/mysql_dump_scheduler/common"
)

// @description:
// SetDatabaseConfig -> set sql.db instance and configure database options
//
//	type DatabaseConfig struct {
//		Host         string   `json:"host"`
//		Port         string   `json:"port"`
//		User         string   `json:"user"`
//		Password     string   `json:"password"`
//		Database     string   `json:"database"`
//		SqlFilesPath SqlFiles `json:"file_path"`     // path to the file to insert the data
//		DumpDirPath  string   `json:"dump_dir_path"` // path to the dump directory`
//	}
func (s *InitializersService) SetDatabaseConfig(db *sql.DB, opts *common.DatabaseConfig) error {
	if opts.Host == "" || opts.Password == "" || opts.User == "" || opts.Database == "" {
		return errors.New("error: some database config oprions is empty")
	}

	if db == nil {
		return errors.New("error: db instance is nil")
	}

	s.dbConfig = opts
	s.db = db
	return nil
}

// initializeDatabaseIfNotExists -> initialize database if not exists
// use only once during the first start of the service
func (s *InitializersService) InitializeDatabaseIfNotExists() error {

	if s.db != nil {
		return errors.New("error: database already initialized")
	}

	connectionStr := s.getConnenctionString()
	connectionRootString := strings.Split(connectionStr, "/")[0] + "/"

	rootDb, err := sql.Open("mysql", connectionRootString)
	if err != nil {
		log.Fatal("InitDatabaseService error: Cannot connect to MySQL server")
	}
	defer rootDb.Close()

	_, err = rootDb.Exec("CREATE DATABASE IF NOT EXISTS " + s.dbConfig.Database)
	if err != nil {
		log.Fatal("InitDatabaseService error: Failed to create database")
	}

	if err := s.runSQLFile(s.dbConfig.SqlFilesPath.TablesFilePath); err != nil {
		return err
	}

	if s.dbConfig.SqlFilesPath.DataFilePath != "" {
		if err := s.runSQLFile(s.dbConfig.SqlFilesPath.DataFilePath); err != nil {
			return err
		}
	}

	log.Println("Database created and initialized with all mock dataschemes successfully")
	return nil
}
