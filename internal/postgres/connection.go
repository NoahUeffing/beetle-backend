package postgres

import (
	"beetle/internal/config"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	gormpq "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func open(conn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Fatalf("Error connecting to database:\n%s\n", err.Error())
	}
	return db
}

func Open(config config.DBsConfig) (*sqlx.DB, *sqlx.DB) {
	return open(config.Read), open(config.Write)
}

func Close(db *sqlx.DB) {
	if db != nil {
		db.Close()
	}
}

func openGorm(conn string) *gorm.DB {
	db, err := gorm.Open(gormpq.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database with gorm:\n%s\n", err.Error())
	}

	return db
}

func OpenGorm(config config.DBsConfig) (*gorm.DB, *gorm.DB) {
	return openGorm(config.Read), openGorm(config.Write)
}
