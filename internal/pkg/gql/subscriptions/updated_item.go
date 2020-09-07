package subscriptions

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// UpdatedItem resolves the updatedItem subscription
func UpdatedItem(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println("Processing subscription UpdatedItem...")

	db := db.FetchConnection()

	rootValue := p.Info.RootValue.(map[string]interface{})
	payload := rootValue["updateItem"].(map[string]interface{})

	tripID := p.Args["tripId"]
	item := &models.Item{}
	if err := db.Where("id = ? AND grocery_trip_id = ?", payload["id"], tripID).First(&item).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return item, nil
}
