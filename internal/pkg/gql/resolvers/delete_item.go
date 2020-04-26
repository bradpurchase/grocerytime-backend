package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// DeleteItemResolver deletes an item by itemId param
func DeleteItemResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	item := &models.Item{}
	if err := db.Where("id = ?", p.Args["itemId"]).First(&item).Error; err != nil {
		return nil, err
	}

	if err := db.Delete(&item).Error; err != nil {
		return nil, err
	}
	return item, err
}
