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

	if s.SchedulerConfig.MakeOpts.RunPath == "" {
		if err := s.dumpScript(); err != nil {
			return nil, fmt.Errorf("failed to create dump file: %v", err)
		}
	}

	cmd := exec.Command("make", "run-dump")
	cmd.Dir = s.SchedulerConfig.MakeOpts.RunPath

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("DB_HOST=%s", s.DatabaseConfig.Host),
		fmt.Sprintf("DB_PORT=%s", s.DatabaseConfig.Port),
		fmt.Sprintf("DB_USER=%s", s.DatabaseConfig.User),
		fmt.Sprintf("DB_PASSWORD=%s", s.DatabaseConfig.Password),
		fmt.Sprintf("DB_NAME=%s", s.DatabaseConfig.Database),
		fmt.Sprintf("DB_DUMP_DIR=%s", s.DatabaseConfig.SqlFilesPath.DumpDirPath),
	)

	output, err := cmd.CombinedOutput()
	log.Printf("createDumpFile output: %s", output)
	if err != nil {
		return nil, fmt.Errorf("failed to create dump file: %v, output: %s", err, output)
	}

	file, err := os.Open(s.DatabaseConfig.SqlFilesPath.DumpDirPath + "/" + s.DatabaseConfig.Database + ".sql")
	if err != nil {
		return nil, fmt.Errorf("failed to open dump file: %v", err)
	}
	defer file.Close()

	return file, nil
}

func (s *SchedulerService) compressDumpFile() error {

	filePath := s.DatabaseConfig.SqlFilesPath.DumpDirPath + "/" + s.DatabaseConfig.Database + ".sql"
	cmd := exec.Command("gzip", "-f", filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compress dump file: %v, output: %s", err, output)
	}

	return nil
}

func (s *SchedulerService) getCompressedFile() (*os.File, error) {

	compressedFilePath := s.DatabaseConfig.SqlFilesPath.DumpDirPath + "/" + s.DatabaseConfig.Database + ".sql.gz"
	file, err := os.Open(compressedFilePath)
	if err != nil {
		log.Printf("getCompressedFile error: failed to open compressed file: %v", err)
		return nil, fmt.Errorf("failed to open compressed file: %v", err)
	}
	return file, nil
}

func (s *SchedulerService) dumpScript() error {

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err := os.MkdirAll(s.DatabaseConfig.SqlFilesPath.DumpDirPath, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	outputFile := filepath.Join(s.DatabaseConfig.SqlFilesPath.DumpDirPath, fmt.Sprintf("%s_%s.sql", s.DatabaseConfig.Database, timestamp))
	cmd := exec.Command("mysqldump",
		"-h", s.DatabaseConfig.Host,
		"-P", s.DatabaseConfig.Port,
		"-u", s.DatabaseConfig.User,
		fmt.Sprintf("-p%s", s.DatabaseConfig.Password),
		s.DatabaseConfig.Database)

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

	s.logChan <- fmt.Sprintf("Database dump created at %s\n", outputFile)

	return nil
}

// @Run -> run the scheduler in following order:
//  1. set the timer
//  2. create the dump file
//  3. compress the dump file
//  4. send the file to the telegram // -> optional
func (s *SchedulerService) run() error {

	file, err := s.createDumpFile()
	if err != nil {
		return err
	}

	s.SchedulerConfig.File = file
	if err := s.compressDumpFile(); err != nil {
		return err
	}

	s.logChan <- "SchedulerService::run: dump file created"

	if s.TelegramConfig.ChatId != "" {
		if err := s.sendFileToTelegram(); err != nil {
			return err
		}
		s.logChan <- "SchedulerService::run: file sent to telegram"
	}

	return nil
}
