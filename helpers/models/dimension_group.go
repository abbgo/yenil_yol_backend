package helpers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/serializations"
)

func GetDimensionsByDimensionGroupID(dimensionGroupID string) (serializations.DimensionGroup, error) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		return serializations.DimensionGroup{}, err
	}
	defer db.Close()

	var dimensionGroup serializations.DimensionGroup
	if err := db.QueryRow(
		context.Background(), `SELECT id,name FROM dimension_groups WHERE id=$1`, dimensionGroupID).
		Scan(&dimensionGroup.ID, &dimensionGroup.Name); err != nil {
		return serializations.DimensionGroup{}, err
	}

	// bu razmer grupba degisli razmerler alynyar
	rowsDimensions, err := db.Query(context.Background(), `SELECT dimension FROM dimensions WHERE dimension_group_id = $1`, dimensionGroup.ID)
	if err != nil {
		return serializations.DimensionGroup{}, err
	}
	defer rowsDimensions.Close()

	for rowsDimensions.Next() {
		var dimension string
		if err := rowsDimensions.Scan(&dimension); err != nil {
			return serializations.DimensionGroup{}, err
		}
		dimensionGroup.Dimensions = append(dimensionGroup.Dimensions, dimension)
	}

	return dimensionGroup, nil
}
