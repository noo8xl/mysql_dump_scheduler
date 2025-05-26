package initializers

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func (s *InitializersService) getConnenctionString() string {
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", s.dbConfig.User, s.dbConfig.Password, s.dbConfig.Host, s.dbConfig.Port, s.dbConfig.Database)
	return str
}

func (s *InitializersService) runSQLFile(filePath string) error {

	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("error: InitDatabaseService failed to read SQL file")
		return err
	}

	_, err = s.db.Exec("USE " + s.dbConfig.Database)
	if err != nil {
		log.Printf("error: InitDatabaseService failed to use database")
		return err
	}

	sqlStatements := strings.Split(string(sqlBytes), ";")
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err = s.db.Exec(stmt)
		if err != nil {
			log.Printf("error: InitDatabaseService failed to execute SQL statement: %v", err)
			return err
		}
	}

	return nil
}
