package postgres_test

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockDBs struct {
	ReadDB    *gorm.DB
	ReadMock  sqlmock.Sqlmock
	WriteDB   *gorm.DB
	WriteMock sqlmock.Sqlmock
}

func open() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		log.Fatal("Error initializing mock database: " + err.Error())
	}

	mock.ExpectPing()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Error creating GORM DB: " + err.Error())
	}

	return gormDB, mock
}

func getTestDBs(t *testing.T) *MockDBs {
	readDB, readMock := open()
	writeDB, writeMock := open()

	t.Cleanup(func() {
		sqlDB, err := readDB.DB()
		if err != nil {
			t.Fatal("Error getting SQL DB from GORM: " + err.Error())
		}
		sqlDB.Close()

		sqlDB, err = writeDB.DB()
		if err != nil {
			t.Fatal("Error getting SQL DB from GORM: " + err.Error())
		}
		sqlDB.Close()
	})

	return &MockDBs{
		ReadMock:  readMock,
		ReadDB:    readDB,
		WriteMock: writeMock,
		WriteDB:   writeDB,
	}
}
