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

func TestRetrieveUserLists_HasLists(t *testing.T) {
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
		AddRow(uuid.NewV4(), userID, "Test List", time.Now(), time.Now()).
		AddRow(uuid.NewV4(), userID, "Test List 2", time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"*").
		WithArgs(userID).
		WillReturnRows(listRows)

	userLists, err := RetrieveUserLists(db, userID)
	require.NoError(t, err)
	assert.Len(t, userLists, 2)
}

//TODO more test cases needed:
// - test that this doesn't return other user's lists
// - test that this includes both lists created and joined
