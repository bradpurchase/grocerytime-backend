package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// ListsResolver returns List records for the current user
//TODO: this needs to return any record where the user has a ListUser association.
// Returning List records will only return lists the user created, not including ones
// they got added to
func ListsResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	user := p.Source.(models.User)
	lists := []models.List{}
	if err := db.Where("user_id = ?", user.ID).Find(&lists).Error; err != nil {
		return nil, err
	}
	return lists, nil
}
