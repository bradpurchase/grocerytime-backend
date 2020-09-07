package subscriptions

import (
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// NewItem resolves the newItem subscription
func NewItem(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println("Processing subscription NewItem...")

	db := db.FetchConnection()

	rootValue := p.Info.RootValue.(map[string]interface{})
	payload := rootValue["addItemToTrip"].(map[string]interface{})
	item := &models.Item{}
	tripID := p.Args["tripId"]
	if err := db.Where("id = ? AND grocery_trip_id = ?", payload["id"], tripID).First(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
