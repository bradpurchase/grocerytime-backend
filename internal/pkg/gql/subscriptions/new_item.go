package subscriptions

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// NewItem resolves the newItem subscription
func NewItem(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println("Processing subscription NewItem...")

	db := db.FetchConnection()
	defer db.Close()

	rootValue := p.Info.RootValue.(map[string]interface{})
	payload := rootValue["addItemToTrip"].(map[string]interface{})

	tripID := p.Args["tripId"]
	item := &models.Item{}
	if err := db.Where("id = ? AND grocery_trip_id = ?", payload["id"], tripID).First(&item).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return item, nil
}
