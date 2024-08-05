package serializations

type ProductComplaint struct {
	ProductID      string `json:"product_id,omitempty"`
	NameTM         string `json:"name_tm,omitempty"`
	NameRU         string `json:"name_ru,omitempty"`
	Image          string `json:"image,omitempty"`
	ComplaintCount int    `json:"complaint_count,omitempty"`
}
