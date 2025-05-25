package initializers

import (
	"database/sql"

	"github.com/noo8xl/mysql_dump_scheduler/common"
)

type initializersService struct {
	dbConfig *common.DatabaseConfig
	db       *sql.DB
}

// @desctiption:
// NewInitializersService -> initialize the service with an empty configs
func NewInitializersService() *initializersService {
	return &initializersService{
		dbConfig: &common.DatabaseConfig{},
		db:       &sql.DB{},
	}
}
