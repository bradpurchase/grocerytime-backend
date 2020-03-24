package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// ListsResolver returns List records for the current user
func ListsResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	// Find the lists this user belongs to
	user := p.Source.(models.User)
	lists := []models.List{}
	query := db.
		Select("lists.*").
		Joins("INNER JOIN list_users ON list_users.list_id = lists.id").
		Where("list_users.user_id = ?", user.ID).
		Find(&lists).
		Error
	if err := query; err != nil {
		return nil, err
	}

	return lists, nil
}
