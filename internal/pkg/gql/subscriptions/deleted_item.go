package subscriptions

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// DeletedItem resolves the deletedItem subscription
func DeletedItem(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println("Processing subscription DeletedItem...")

	db := db.FetchConnection()
	defer db.Close()

	rootValue := p.Info.RootValue.(map[string]interface{})
	payload := rootValue["deleteItem"].(map[string]interface{})

	tripID := p.Args["tripId"]
	item := &models.Item{}
	if err := db.Where("id = ? AND grocery_trip_id = ?", payload["id"], tripID).First(&item).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return item, err // note: we return item here so that the subscriber gets _some_ result
	}

	return item, nil
}
