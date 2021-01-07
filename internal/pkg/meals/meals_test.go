package meals

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
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

func (s *Suite) TestRetrieveMeals_InvalidWeekNumber() {
	userID := uuid.NewV4()
	weekNumber := 54
	year := 2021
	_, e := RetrieveMeals(userID, weekNumber, year)
	require.Error(s.T(), e)
	assert.Contains(s.T(), e.Error(), "number of week can't be less than 1 or greater than 53")
}

func (s *Suite) TestRetrieveMeals_NoMeals() {
	userID := uuid.NewV4()
	weekNumber := 1
	year := 2021
	s.mock.ExpectQuery("^SELECT meals.* FROM \"meals\"*").
		WithArgs(userID, "2021-01-04", "2021-01-10").
		WillReturnRows(sqlmock.NewRows([]string{}))

	meals, err := RetrieveMeals(userID, weekNumber, year)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(meals), 0)
}

func (s *Suite) TestRetrieveMeals_ExistingMeals() {
	userID := uuid.NewV4()
	weekNumber := 1
	year := 2021

	meal1ID := uuid.NewV4()
	meal2ID := uuid.NewV4()
	mealRows := sqlmock.NewRows([]string{"id"}).
		AddRow(meal1ID).
		AddRow(meal2ID)
	s.mock.ExpectQuery("^SELECT meals.* FROM \"meals\"*").
		WithArgs(userID, "2021-01-04", "2021-01-10").
		WillReturnRows(mealRows)

	mealUserRows := sqlmock.NewRows([]string{"id", "meal_id", "user_id"}).
		AddRow(uuid.NewV4(), meal1ID, userID).
		AddRow(uuid.NewV4(), meal2ID, userID)
	s.mock.ExpectQuery("^SELECT (.+) FROM \"meal_users\"*").
		WithArgs(meal1ID, meal2ID).
		WillReturnRows(mealUserRows)

	meals, err := RetrieveMeals(userID, weekNumber, year)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), len(meals), 2)
	assert.Equal(s.T(), meals[0].ID, meal1ID)
	assert.Equal(s.T(), meals[1].ID, meal2ID)
}

func (s *Suite) TestRetrieveMeal_NotFound() {
	mealID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT meals.* FROM \"meals\"*").
		WithArgs(mealID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := RetrieveMealForUser(mealID, userID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "record not found")
}

func (s *Suite) TestRetrieveMeal_Found() {
	mealID := uuid.NewV4()
	userID := uuid.NewV4()
	s.mock.ExpectQuery("^SELECT meals.* FROM \"meals\"*").
		WithArgs(mealID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mealID))

	meal, err := RetrieveMealForUser(mealID, userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), meal.ID, mealID)
}
