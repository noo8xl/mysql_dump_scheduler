package initializers_test

import (
	"database/sql"
	"testing"

	"github.com/noo8xl/mysql_dump_scheduler/common"
	"github.com/noo8xl/mysql_dump_scheduler/initializers"
)

var (
	initSvc  *initializers.InitializersService
	db       *sql.DB
	dbConfig = common.DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "password",
		Database: "test",
		SqlFilesPath: &common.SqlFiles{
			TablesFilePath: "/home/noo8xl/Documents/workspace/ANVIL/anvil-reviews/database/schemes/db.sql",
			DataFilePath:   "/home/noo8xl/Documents/workspace/ANVIL/anvil-reviews/database/data/data.sql",
			DumpDirPath:    "/home/noo8xl/Documents/workspace/ANVIL/anvil-reviews/database/backups",
		},
	}
)

func TestInitializeService(t *testing.T) {
	initSvc = initializers.NewInitializersService()
	if err := initSvc.SetDatabaseConfig(db, dbConfig); err != nil {
		t.Fatalf("error: failed to set database config: %v", err)
	}
	if err := initSvc.InitializeDatabaseIfNotExists(); err != nil {
		t.Fatalf("error: failed to initialize database: %v", err)
	}
}
