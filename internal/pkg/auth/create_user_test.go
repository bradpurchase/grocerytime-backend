package auth

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateUser_InvalidPassword(t *testing.T) {
	dbMock, _, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	email := "test@example.com"
	_, e := CreateUser(db, email, "")
	require.Error(t, e)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	email := "test@example.com"
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(email))

	_, e := CreateUser(db, email, "password")
	require.Error(t, e)
	assert.Equal(t, e.Error(), "An account with this email address already exists")
}

func TestCreateUser_UserCreated(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	require.NoError(t, err)

	email := "test@example.com"
	storeName := "My Grocery Store"

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userID := uuid.NewV4()
	mock.ExpectQuery("^INSERT INTO \"users\" (.+)$").
		WithArgs(email, sqlmock.AnyArg(), "", "", AnyTime{}, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	clientID := uuid.NewV4()
	mock.ExpectQuery("^SELECT \"id\" FROM \"api_clients\"*").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(clientID))
	mock.ExpectQuery("INSERT INTO \"auth_tokens\" (\"access_token\",\"refresh_token\",\"expires_in\",\"created_at\",\"updated_at\",\"client_id\",\"user_id\")*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), AnyTime{}, AnyTime{}, AnyTime{}, clientID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "client_id", "user_id"}).AddRow(uuid.NewV4(), clientID, userID))

	storeID := uuid.NewV4()
	mock.ExpectQuery("^INSERT INTO \"stores\" (.+)$").
		WithArgs(storeName, AnyTime{}, AnyTime{}, nil, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(storeID, userID))
	mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"store_id"}).AddRow(storeID))

	categories := fetchCategories()
	for i := range categories {
		mock.ExpectQuery("^INSERT INTO \"store_categories\" (.+)$").
			WithArgs(storeID, categories[i], AnyTime{}, AnyTime{}, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	}

	mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	user, err := CreateUser(db, email, "password")
	require.NoError(t, err)
	assert.Equal(t, user.Email, email)
}

// TODO: duplicated code with the store model... DRY this up
func fetchCategories() [20]string {
	categories := [20]string{
		"Produce",
		"Bakery",
		"Meat",
		"Seafood",
		"Dairy",
		"Cereal",
		"Baking",
		"Dry Goods",
		"Canned Goods",
		"Frozen Foods",
		"Cleaning",
		"Paper Products",
		"Beverages",
		"Candy & Snacks",
		"Condiments",
		"Personal Care",
		"Baby",
		"Alcohol",
		"Pharmacy",
		"Misc.",
	}
	return categories
}
