package serializations

type ShopOwner struct {
	ID          string `json:"id,omitempty"`
	FullName    string `json:"full_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}
