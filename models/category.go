package models

import (
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"

	"gopkg.in/guregu/null.v4"
)

type Category struct {
	ID               string      `json:"id,omitempty"`
	NameTM           string      `json:"name_tm,omitempty" binding:"required"`
	NameRU           string      `json:"name_ru,omitempty" binding:"required"`
	Image            string      `json:"image,omitempty"`
	SlugTM           string      `json:"slug_tm,omitempty"`
	SlugRU           string      `json:"slug_ru,omitempty"`
	DimensionGroupID string      `json:"dimension_group_id,omitempty" binding:"required"`
	ParentCategoryID null.String `json:"parent_category_id,omitempty"`
}

func ValidateCategory(category Category, isCreateFunction bool) error {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := helpers.ValidateRecordByID("dimension_groups", category.DimensionGroupID, "NULL", db); err != nil {
		return err
	}

	if !isCreateFunction {
		if category.ID == "" {
			return errors.New("category_id is required")
		}

		if err := helpers.ValidateRecordByID("categories", category.ID, "NULL", db); err != nil {
			return errors.New("record not found")
		}
	}

	// validate parentCategoryID
	if category.ParentCategoryID.String != "" {
		if isCreateFunction {
			if category.Image != "" {
				return errors.New("child cannot be an image of the category")
			}
		}

		if err := helpers.ValidateRecordByID("categories", category.ParentCategoryID.String, "NULL", db); err != nil {
			return err
		}

		return nil
	} else {
		if isCreateFunction {
			if category.Image == "" {
				return errors.New("parent category image is required")
			}
		}
	}

	return nil
}
