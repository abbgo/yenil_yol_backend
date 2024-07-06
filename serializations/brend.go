package serializations

import "github/abbgo/yenil_yol/backend/helpers"

type BrendQuery struct {
	helpers.StandartQuery
	Search string `form:"search"`
}
