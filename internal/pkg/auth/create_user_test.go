package auth

import (
	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestCreateUser_InvalidPassword() {
	email := "test@example.com"
	_, e := CreateUser(email, "", uuid.NewV4())
	require.Error(s.T(), e)
}

func (s *Suite) TestCreateUser_DuplicateEmail() {
	email := "test@example.com"
	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(email))

	_, e := CreateUser(email, "password", uuid.NewV4())
	require.Error(s.T(), e)
	assert.Equal(s.T(), e.Error(), "An account with this email address already exists")
}

func (s *Suite) TestCreateUser_UserCreated() {
	email := "test@example.com"
	storeName := "My Grocery Store"

	s.mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"users\" (.+)$").
		WithArgs(email, sqlmock.AnyArg(), "", "", AnyTime{}, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	clientID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"auth_tokens\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), AnyTime{}, AnyTime{}, AnyTime{}, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "client_id", "user_id"}).AddRow(uuid.NewV4(), clientID, userID))

	storeID := uuid.NewV4()
	s.mock.ExpectQuery("^INSERT INTO \"stores\" (.+)$").
		WithArgs(storeName, AnyTime{}, AnyTime{}, nil, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(storeID, userID))
	s.mock.ExpectQuery("^INSERT INTO \"store_users\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "", true, true, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	categories := fetchCategories()
	for i := range categories {
		s.mock.ExpectQuery("^INSERT INTO \"store_categories\" (.+)$").
			WithArgs(sqlmock.AnyArg(), categories[i], AnyTime{}, AnyTime{}, nil).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))
	}

	s.mock.ExpectQuery("^INSERT INTO \"grocery_trips\" (.+)$").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), false, false, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

	user, err := CreateUser(email, "password", clientID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), user.Email, email)
	assert.NotNil(s.T(), user.Tokens)
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
