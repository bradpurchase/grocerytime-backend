package list

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveUserLists_NoLists(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	userID := uuid.NewV4()
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	userLists, err := RetrieveUserLists(db, userID)
	require.NoError(t, err)
	assert.Equal(t, userLists, []models.List{})
}

func TestRetrieveUserLists_HasListsCreated(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	userID := uuid.NewV4()
	listRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), userID, "Provigo Grocery List", time.Now(), time.Now()).
		AddRow(uuid.NewV4(), userID, "Beer Store List", time.Now(), time.Now())
	mock.
		ExpectQuery("^SELECT lists.* FROM \"lists\"*").
		WithArgs(userID).
		WillReturnRows(listRows)

	db.Exec("INSERT INTO lists (name, user_id) VALUES (?, ?)", "Not My List", uuid.NewV4())

	userLists, err := RetrieveUserLists(db, userID)
	require.NoError(t, err)
	assert.Len(t, userLists, 2)
	assert.Equal(t, userLists.([]models.List)[0].Name, "Provigo Grocery List")
	assert.Equal(t, userLists.([]models.List)[0].UserID, userID)
	assert.Equal(t, userLists.([]models.List)[1].Name, "Beer Store List")
	assert.Equal(t, userLists.([]models.List)[1].UserID, userID)
}

func TestRetrieveUserLists_HasListsCreatedAndJoined(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open("postgres", dbMock)
	require.NoError(t, err)

	userID := uuid.NewV4()
	sharingUserID := uuid.NewV4()
	sharedListID := uuid.NewV4()
	listRows := sqlmock.
		NewRows([]string{
			"id",
			"user_id",
			"name",
			"created_at",
			"updated_at",
		}).
		AddRow(uuid.NewV4(), userID, "Provigo Grocery List", time.Now(), time.Now()).
		AddRow(uuid.NewV4(), userID, "Beer Store List", time.Now(), time.Now()).
		AddRow(sharedListID, sharingUserID, "Shared List", time.Now(), time.Now())
	mock.
		ExpectQuery("^SELECT lists.* FROM \"lists\"*").
		WithArgs(userID).
		WillReturnRows(listRows)

	db.Exec("INSERT INTO list_users (list_id, user_id) VALUES (?, ?)", sharedListID, userID)

	userLists, err := RetrieveUserLists(db, userID)
	require.NoError(t, err)
	assert.Len(t, userLists, 3)
	assert.Equal(t, userLists.([]models.List)[0].Name, "Provigo Grocery List")
	assert.Equal(t, userLists.([]models.List)[0].UserID, userID)
	assert.Equal(t, userLists.([]models.List)[1].Name, "Beer Store List")
	assert.Equal(t, userLists.([]models.List)[2].Name, "Shared List")
	assert.Equal(t, userLists.([]models.List)[2].UserID, sharingUserID)
}
