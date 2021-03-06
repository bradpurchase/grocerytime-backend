package stores

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestCreateStore_DupeStore() {
	userID := uuid.NewV4()
	storeName := "My Dupe Store"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(uuid.NewV4(), userID, storeName))

	_, e := CreateStore(userID, storeName)
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "you already added a store with this name")
}

func (s *Suite) TestCreateStore_Created() {
	userID := uuid.NewV4()
	storeName := "My New Store"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"stores\"*").
		WithArgs(storeName, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	storeID := uuid.NewV4()
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("^INSERT INTO \"stores\" (.+)$").
		WithArgs(storeName, sqlmock.AnyArg(), AnyTime{}, AnyTime{}, nil, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(storeID, userID))

	storeUserID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storeUserID))
	s.mock.ExpectQuery("^INSERT INTO \"store_user_preferences\" (.+)$").
		WithArgs(storeUserID, false, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	categories := fetchCategories()
	for i := range categories {
		s.mock.ExpectQuery("^INSERT INTO \"store_categories\" (.+)$").
			WithArgs(sqlmock.AnyArg(), categories[i], AnyTime{}, AnyTime{}, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	}

	currentTime := time.Now()
	tripName := currentTime.Format("Jan 2, 2006")
	likeTripName := fmt.Sprintf("%%%s%%", tripName)
	s.mock.ExpectQuery("^SELECT count*").
		WithArgs(likeTripName, sqlmock.AnyArg()).
		WillReturnRows(s.mock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(sqlmock.AnyArg(), tripName, false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	s.mock.ExpectCommit()

	store, err := CreateStore(userID, storeName)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), store.Name, storeName)
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
