package resolvers

import (
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/graphql-go/graphql"
)

// CreateRecipeResolver creates a new recipe
func CreateRecipeResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	fmt.Println("args", p.Args)
	return nil, nil
}
