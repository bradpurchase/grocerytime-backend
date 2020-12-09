package trips

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
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

func (s *Suite) TestRetrieveCurrentStoreTripForUser_UserNotMemberOfStore() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4()}
	_, err := RetrieveCurrentStoreTripForUser(storeID, user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "user is not a member of this store")
}

func (s *Suite) TestRetrieveCurrentStoreTripForUser_TripNotAssociatedWithStore() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	_, err := RetrieveCurrentStoreTripForUser(storeID, user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "could not find trip associated with this store")
}

func (s *Suite) TestRetrieveCurrentStoreTripForUser_FoundResult() {
	storeID := uuid.NewV4()
	user := models.User{ID: uuid.NewV4(), Email: "test@example.com"}
	s.mock.ExpectQuery("^SELECT (.+) FROM \"store_users\"*").
		WithArgs(storeID, user.ID, user.Email).
		WillReturnRows(s.mock.NewRows([]string{"id", "user_id"}).AddRow(uuid.NewV4(), user.ID))

	tripID := uuid.NewV4()
	tripName := "Week of May 31"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(s.mock.NewRows([]string{"id", "name"}).AddRow(tripID, tripName))

	trip, err := RetrieveCurrentStoreTripForUser(storeID, user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), trip.Name, tripName)
}

func (s *Suite) TestRetrieveTrips_UserNotActive() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(storeID, userID, true).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))

	_, e := RetrieveTrips(storeID, userID, false)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "user is not active in this store")
}

func (s *Suite) TestRetrieveTrips_Found() {
	storeID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(storeID, userID, true).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(1))

	rows := s.mock.NewRows([]string{"id"}).
		AddRow(uuid.NewV4()).
		AddRow(uuid.NewV4())
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(storeID, false).
		WillReturnRows(rows)

	trips, err := RetrieveTrips(storeID, userID, false)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(trips), 2)
}

func (s *Suite) TestRetrieveTrip_NotFound() {
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{}))

	_, e := RetrieveTrip(tripID)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "trip not found")
}

func (s *Suite) TestRetrieveTrip_Found() {
	tripID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT (.+) FROM \"grocery_trips\"*").
		WithArgs(tripID).
		WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(tripID))

	trip, err := RetrieveTrip(tripID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), trip.ID, tripID)
}
