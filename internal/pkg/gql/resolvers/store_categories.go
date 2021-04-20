package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// StoreCategoriesResolver resolves the storeCategories query
func StoreCategoriesResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	storeID := p.Args["storeId"]
	var storeCategories []models.StoreCategory
	query := db.Manager.
		Joins("INNER JOIN store_users ON store_users.store_id = store_categories.store_id").
		Where("store_categories.store_id = ?", storeID).
		Where("store_users.user_id = ?", user.ID).
		Find(&storeCategories).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return storeCategories, nil
}
