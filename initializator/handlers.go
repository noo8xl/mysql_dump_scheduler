package initializator

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

// initializeDatabaseIfNotExists -> initialize database if not exists
// use only once during the first start of the service
func (s *initializatorService) InitializeDatabaseIfNotExists() error {

	// s.db = s.connectToDatabase()
	// defer s.db.Close()

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
