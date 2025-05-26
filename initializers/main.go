package initializers

import (
	"database/sql"

	"github.com/noo8xl/mysql_dump_scheduler/common"
)

type InitializersService struct {
	dbConfig *common.DatabaseConfig
	db       *sql.DB
}

// @desctiption:
// NewInitializersService -> initialize the service with an empty configs
func NewInitializersService() *InitializersService {
	return &InitializersService{
		dbConfig: &common.DatabaseConfig{},
		db:       &sql.DB{},
	}
}
