# mysql_dump_scheduler


## Usage

Just a simple script to dump the database and compress it. It also can be sent to the telegram (optional).

set the options before the initialization:

SchedulerOpts {
  DatabaseConfig {
    Host: "localhost",
    Port: "3306",
    User: "root",
    Password: "password",
    Name: "database_name",
  }
  TelegramConfig {
    ChatId: "1234567890",
    Token: "your-telegram_bot_token", 
  }
  SchedulerConfig {
    Duration: 1 * time.Hour, // set a duration to run the scheduler
    Path: "./dumps", // path to the dir where the <backup> folder will be created to store the dump file
  }
}

to Init the service -> InitScheduler(opts SchedulerOpts)

to run the service -> Bootstrap()

to send the file to the telegram -> SendFile() (optional)


