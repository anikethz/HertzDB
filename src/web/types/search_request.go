package hertzTypes

type Field struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
type SearchRequest struct {
	Field Field `json:"field"`
}
