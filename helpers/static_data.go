package helpers

var CreatedStatuses map[string]int8 = map[string]int8{
	"wait":     0,
	"rejected": 1,
	"success":  2,
}

var Genders map[string]int8 = map[string]int8{
	"male":   0,
	"female": 1,
	"child":  2,
}
