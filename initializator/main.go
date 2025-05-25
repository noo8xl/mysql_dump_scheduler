package initializator

import (
	"database/sql"
	"errors"

	"github.com/noo8xl/mysql_dump_scheduler/common"
)

type initializatorService struct {
	dbConfig *common.DatabaseConfig
	db       *sql.DB
}

func NewInitializatorService() *initializatorService {
	return &initializatorService{
		dbConfig: &common.DatabaseConfig{},
		db:       &sql.DB{},
	}
}

func (s *initializatorService) SetDatabaseConfig(db *sql.DB, opts common.DatabaseConfig) error {
	if opts.Host == "" || opts.Password == "" || opts.User == "" || opts.Database == "" {
		return errors.New("error: some database config oprions is empty")
	}

	if db == nil {
		return errors.New("error: db instance is nil")
	}

	s.dbConfig = &opts
	s.db = db
	return nil
}
