package categories

type CategoryPage struct {
	Categories []CategoryView `json:"categories"`
	Total      int64          `json:"total"`
}

type CategoryView struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
