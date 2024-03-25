package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
)

type Dimension struct {
	ID               string `json:"id,omitempty"`
	Dimension        string `json:"dimension" binding:"required"`
	DimensionGroupID string `json:"dimension_group_id,omitempty" binding:"required"`
}

type DimensionQuery struct {
	IsDeleted        bool   `form:"is_deleted"`
	DimensionGroupID string `form:"dimension_group_id" validate:"required"`
}

func ValidateDimension(dimension Dimension, forUpdate bool) error {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := helpers.ValidateRecordByID("dimension_groups", dimension.DimensionGroupID, "NULL", db); err != nil {
		return err
	}

	if forUpdate {
		if dimension.ID == "" {
			return errors.New("id is required")
		}

		if err := helpers.ValidateRecordByID("dimensions", dimension.ID, "NULL", db); err != nil {
			return err
		}

		var id string
		db.QueryRow(context.Background(), "SELECT id FROM dimensions WHERE dimension=$1 AND deleted_at IS NULL", dimension.Dimension).Scan(&id)
		if dimension.ID != id && id != "" {
			return errors.New("this dimension already exists")
		}
	} else {
		var id string
		db.QueryRow(context.Background(), "SELECT id FROM dimensions WHERE dimension=$1 AND deleted_at IS NULL", dimension.Dimension).Scan(&id)
		if id != "" {
			return errors.New("this dimension already exists")
		}
	}

	return nil
}
