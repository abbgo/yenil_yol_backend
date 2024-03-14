package models

import (
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
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
	ID               string `json:"id,omitempty" binding:"required"`
	NameTM           string `json:"name_tm,omitempty" binding:"required"`
	NameRU           string `json:"name_ru,omitempty" binding:"required"`
	Image            string `json:"image,omitempty"`
	DimensionGroupID string `json:"dimension_group_id,omitempty" binding:"required"`
}

func ValidateCategory(dimensionGroupID string) error {
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
