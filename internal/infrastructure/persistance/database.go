package persistance

import (
	"database/sql"
	"log"
	"math-spark/internal/infrastructure/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
    DB *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
    log.Printf("Attempting to connect to MySQL with URL: %s", cfg.DatabaseURL)
    db, err := sql.Open("mysql", cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }

    for i := 0; i < 10; i++ {
        if err := db.Ping(); err == nil {
            break
        }
        log.Printf("Failed to ping MySQL, retrying in 1s: %v", err)
        time.Sleep(time.Second)
    }
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping MySQL after retries: %v", err)
    }

    driver, err := mysql.WithInstance(db, &mysql.Config{})
    if err != nil {
        log.Printf("Failed to create migration driver: %v", err)
        return &Database{DB: db}, nil
    }
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "mysql",
        driver,
    )
    if err != nil {
        log.Printf("Failed to initialize migrations: %v", err)
        return &Database{DB: db}, nil
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Printf("Failed to apply migrations, continuing without migrations: %v", err)
    } else {
        log.Println("Successfully applied migrations")
    }

    log.Println("Successfully connected to MySQL")
    return &Database{DB: db}, nil
}

func (d *Database) Close() error {
    return d.DB.Close()
}