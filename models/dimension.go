package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
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

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var dimension_group_id string
	db.QueryRow(context.Background(), "SELECT id FROM dimension_groups WHERE id = $1 AND deleted_at IS NULL", dimensionGroupID).Scan(&dimension_group_id)

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if dimensionGroupID == "" {
		return errors.New("dimension group not found")
	}

	return nil
}
