package subscriptions

import (
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// NewItemInTrip resolves the newItemInList subscription
func NewItemInTrip(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println("Processing subscription NewItemInTrip...")

	db := db.FetchConnection()
	defer db.Close()

	rootValue := p.Info.RootValue.(map[string]interface{})
	payload := rootValue["addItemToTrip"].(map[string]interface{})

	tripID := p.Args["tripId"]
	item := &models.Item{}
	if err := db.Where("id = ? AND trip_id = ?", payload["id"], tripID).First(&item).Error; err != nil {
		return nil, err
	}

	return item, nil
}
