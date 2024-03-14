package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
)

type Category struct {
	ID               string `json:"id,omitempty"`
	NameTM           string `json:"name_tm,omitempty" binding:"required"`
	NameRU           string `json:"name_ru,omitempty" binding:"required"`
	Image            string `json:"image,omitempty"`
	SlugTM           string `json:"slug_tm,omitempty"`
	SlugRU           string `json:"slug_ru,omitempty"`
	DimensionGroupID string `json:"dimension_group_id,omitempty" binding:"required"`
	CreatedAt        string `json:"-"`
	UpdatedAt        string `json:"-"`
	DeletedAt        string `json:"-"`
}

type CategoryUpdate struct {
	ID     string `json:"id,omitempty" binding:"required"`
	NameTM string `json:"name_tm,omitempty" binding:"required"`
	NameRU string `json:"name_ru,omitempty" binding:"required"`
	Image  string `json:"image,omitempty"`
}

func ValidateCategory(dimensionGroupID string) error {
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
