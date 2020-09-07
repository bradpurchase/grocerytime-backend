package db

import (
	"testing"

	"gorm.io/gorm"
	// Import postgres dialect for gorm
	_ "gorm.io/gorm/dialects/postgres"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFetchConnection(t *testing.T) {
	dbmock, _, _ := sqlmock.New()
	_, err := gorm.Open("postgres", dbmock)
	if err != nil {
		t.Fatalf("main_test.go: TestDBConnection error %v", err)
	}
}

// func TestFactory(t *testing.T) {
// 	db := TestFetchConnection(t)
// }
