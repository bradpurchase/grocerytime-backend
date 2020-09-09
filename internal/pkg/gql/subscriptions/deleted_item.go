package subscriptions

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// DeletedItem resolves the deletedItem subscription
func DeletedItem(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println("Processing subscription DeletedItem...")

	rootValue := p.Info.RootValue.(map[string]interface{})
	payload := rootValue["deleteItem"].(map[string]interface{})

	tripID := p.Args["tripId"]
	item := &models.Item{}
	query := db.Manager.
		Where("id = ? AND grocery_trip_id = ?", payload["id"], tripID).
		First(&item).
		Error
	if err := query; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return item, err // note: we return item here so that the subscriber gets _some_ result
	}

	return item, nil
}
