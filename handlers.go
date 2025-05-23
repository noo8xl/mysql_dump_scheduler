package scheduler

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// @createDumpFile -> create a dump file with the given database config
// @return *os.File if success, error otherwise
func (s *SchedulerService) createDumpFile() (*os.File, error) {

	dumpDir := s.getDumpDir()
	// runDir := fmt.Sprintf("%s/../../", dumpDir)

	// cmd := exec.Command("make", "run-dump")
	// cmd.Dir = runDir

	// log.Printf("createDumpFile cmd.Dir: %s", cmd.Dir)

	// cmd.Env = append(os.Environ(),
	// 	fmt.Sprintf("DB_HOST=%s", s.DatabaseConfig.Host),
	// 	fmt.Sprintf("DB_PORT=%s", s.DatabaseConfig.Port),
	// 	fmt.Sprintf("DB_USER=%s", s.DatabaseConfig.User),
	// 	fmt.Sprintf("DB_PASSWORD=%s", s.DatabaseConfig.Password),
	// 	fmt.Sprintf("DB_NAME=%s", s.DatabaseConfig.Name),
	// 	fmt.Sprintf("DB_DUMP_DIR=%s", dumpDir),
	// )

	// output, err := cmd.CombinedOutput()
	// log.Printf("createDumpFile output: %s", output)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create dump file: %v, output: %s", err, output)
	// }

	if err := s.dumpScript(); err != nil {
		return nil, fmt.Errorf("failed to create dump file: %v", err)
	}

	file, err := os.Open(dumpDir + s.DatabaseConfig.Name + ".sql")
	if err != nil {
		return nil, fmt.Errorf("failed to open dump file: %v", err)
	}
	defer file.Close()

	return file, nil
}

func (s *SchedulerService) compressDumpFile() error {

	filePath := s.SchedulerConfig.Path + "/" + s.DatabaseConfig.Name + ".sql"
	cmd := exec.Command("gzip", "-f", filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compress dump file: %v, output: %s", err, output)
	}

	return nil
}

func (s *SchedulerService) getCompressedFile() (*os.File, error) {

	compressedFilePath := s.SchedulerConfig.Path + "/" + s.DatabaseConfig.Name + ".gz"
	file, err := os.Open(compressedFilePath)
	if err != nil {
		log.Printf("getCompressedFile error: failed to open compressed file: %v", err)
		return nil, fmt.Errorf("failed to open compressed file: %v", err)
	}
	return file, nil
}

func (s *SchedulerService) getDumpDir() string {
	return fmt.Sprintf("%s/backups", s.SchedulerConfig.File.Name())
}

func (s *SchedulerService) dumpScript() error {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbDumpDir := os.Getenv("DB_DUMP_DIR")

	timestamp := time.Now().Format("20060102150405")

	if err := os.MkdirAll(dbDumpDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	outputFile := filepath.Join(dbDumpDir, fmt.Sprintf("%s_%s.sql", dbName, timestamp))
	cmd := exec.Command("mysqldump",
		"-h", dbHost,
		"-P", dbPort,
		"-u", dbUser,
		fmt.Sprintf("-p%s", dbPassword),
		dbName)

	outfile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running mysqldump: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database dump created at %s\n", outputFile)

	return nil
}
