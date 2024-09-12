package serializations

type DimensionGroup struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Dimensions []string `json:"dimensions"`
}
