package models

import (
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
)

type Dimension struct {
	ID               string `json:"id,omitempty"`
	Dimension        string `json:"dimension" binding:"required"`
	DimensionGroupID string `json:"dimension_group_id,omitempty" binding:"required"`
	CreatedAt        string `json:"-"`
	UpdatedAt        string `json:"-"`
	DeletedAt        string `json:"-"`
}

func ValidateDimension(dimensionGroupID string) error {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := helpers.ValidateRecordByID("dimension_groups", dimensionGroupID, "NULL", db); err != nil {
		return err
	}

	return nil
}
