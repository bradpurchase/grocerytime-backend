package user

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Suite struct {
	suite.Suite

	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		dbMock *sql.DB
		err    error
	)

	dbMock, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(s.T(), err)

	db.Manager = s.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
