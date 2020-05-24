package subscriptions

import (
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// NewItemInList resolves the newItemInList subscription
func NewItemInList(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println("Processing subscription NewItemInList...")

	db := db.FetchConnection()
	defer db.Close()

	rootValue := p.Info.RootValue.(map[string]interface{})
	payload := rootValue["addItemToList"].(map[string]interface{})

	listID := p.Args["list_id"]
	item := &models.Item{}
	if err := db.Where("id = ? AND list_id = ?", payload["id"], listID).First(&item).Error; err != nil {
		return nil, err
	}

	return item, nil
}
