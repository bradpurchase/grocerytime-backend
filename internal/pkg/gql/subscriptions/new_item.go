package subscriptions

import (
	"fmt"
	"log"

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
	log.Println(rootValue)
	payload := rootValue["addItemToTrip"].(map[string]interface{})
	log.Println(payload)

	tripID := p.Args["tripId"]
	log.Println(tripID)
	item := &models.Item{}
	if err := db.Where("id = ? AND grocery_trip_id = ?", payload["id"], tripID).First(&item).Error; err != nil {
		return nil, err
	}

	return item, nil
}
