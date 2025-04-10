// TODO: Implement
package postgres_test

import (
	"beetle/internal/postgres"
	"errors"
	"testing"

	"github.com/matryer/is"
)

func getHealthCheckService(dbs *MockDBs) *postgres.HealthCheckService {
	return &postgres.HealthCheckService{
		Read:  dbs.ReadDB,
		Write: dbs.WriteDB,
	}
}

func TestHealthCheckService_Check_Success(t *testing.T) {
	dbs := getTestDBs(t)
	is := is.New(t)
	hcs := getHealthCheckService(dbs)
	dbs.ReadMock.ExpectPing()
	dbs.WriteMock.ExpectPing()
	status := hcs.Check()
	is.Equal(status.Up, true)                         // reported service down
	is.Equal(status.Name, postgres.HEALTH_CHECK_NAME) // incorrect name
	is.Equal(status.Messages, []string{})
}

func TestHealthCheckService_Check_ReadFailure(t *testing.T) {
	dbs := getTestDBs(t)
	is := is.New(t)
	hcs := getHealthCheckService(dbs)
	err := errors.New("database unavailable")
	dbs.ReadMock.ExpectPing().WillReturnError(err)
	dbs.WriteMock.ExpectPing()
	status := hcs.Check()
	is.Equal(status.Up, false)                        // reported service down
	is.Equal(status.Name, postgres.HEALTH_CHECK_NAME) // incorrect name
	is.Equal(status.Messages, []string{"Read DB Unavailable: " + err.Error()})
}

func TestHealthCheckService_Check_WriteFailure(t *testing.T) {
	dbs := getTestDBs(t)
	is := is.New(t)
	hcs := getHealthCheckService(dbs)
	err := errors.New("database unavailable")
	dbs.ReadMock.ExpectPing()
	dbs.WriteMock.ExpectPing().WillReturnError(err)
	status := hcs.Check()
	is.Equal(status.Up, false)                        // reported service down
	is.Equal(status.Name, postgres.HEALTH_CHECK_NAME) // incorrect name
	is.Equal(status.Messages, []string{"Write DB Unavailable: " + err.Error()})
}
