package catalogapi

type ProductPage struct {
	Products []ProductView `json:"products"`
	Total    int64         `json:"total"`
}

type ProductView struct {
	Category string  `json:"category"`
	Code     string  `json:"code"`
	Price    float64 `json:"price"`

	Variants []VariantView `json:"variants,omitempty"`
}

type VariantView struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}
