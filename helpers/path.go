package helpers

import "os"

var ServerPath = os.Getenv("SERVER_PATH")

type StandartQuery struct {
	IsDeleted bool `form:"is_deleted"`
	Limit     int  `form:"limit" validate:"required,min=10,max=20"`
	Page      int  `form:"page" validate:"required,min=1"`
}
