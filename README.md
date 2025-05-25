# mysql_dump_scheduler


## Usage

It's just a simple lib to interact with mysql database.

The main functional possibilities are: 
- create a mysql dump file;
- compress it;
- send it via telegram (optional);
- init database if not exists;
- set default db data;


### Example

```go
import (
  mds "github.com/noo8xl/mysql_dump_scheduler"
)

func main( ) {

  var db *sql.DB // call your db instance here * 

  // db params
  opts := mds.DatabaseConfig {
    Host: "localhost",
    Port: "3306",
    User: "root",
    Password: "password",
    Database: "database_name",
    SqlFilesPath: &SqlFiles{
      TablesFilePath: "path/to-your/tables.sql",
      DataFilePath: "path/to-your/data.sql",
    },
    DumpDirPath: "path/to-your/backup/dir" // optional (use for the scheduler)
  }

  // telegram opts ( optional, use only if you want to send the file to your telegram bot )
  tgOpts := mds.TelegramConfig {
    ChatId: "your-chat-id",
    Token: "your-bot-token",
  }

  // scheduler config 
  schedulerOpts := mds.SchedulerConfig {
    Path: "path/to/the/dump.sql",
    Duration: 24 * time.Hour,
    MakeOpts: &MakeOpts{ // can be omitted. use if u want to use a makefile options
      RunPath: "path/to/your/Makefile",
    }      
  }

  // call the service to init your database and set default data (optional)
  initSvc := mds.InitService()

  initSvc.SetDatabaseConfig(db, opts)
  initSvc.InitializeDatabaseIfNotExists()

  // to run scheduler 
  schedulerSvc := mds.InitScheduler()

  schedulerSvc.SetDatabaseConfig(opts)
  schedulerSvc.SetTelegramConfig(tgOpts)
  schedulerSvc.SetSchedulerConfig(schedulerOpts)

  
  if err := schedulerSvc.Bootstrap(context.Background()); err != nil {
    log.Fatalf("error: scheduler fail with err: %v", err)
  }
}

```go
