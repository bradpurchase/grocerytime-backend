package db

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFetchConnection(t *testing.T) {
	dbMock, _, _ := sqlmock.New()
	_, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})

	if err != nil {
		t.Fatalf("main_test.go: TestDBConnection error %v", err)
	}
}

// func TestFactory(t *testing.T) {
// 	db := TestFetchConnection(t)
// }
