package helpers

import "os"

var ServerPath = os.Getenv("STATIC_FILE_PATH")

type StandartQuery struct {
	IsDeleted bool `form:"is_deleted"`
	Limit     int  `form:"limit" validate:"required,min=1,max=100"`
	Page      int  `form:"page" validate:"required,min=1"`
}
