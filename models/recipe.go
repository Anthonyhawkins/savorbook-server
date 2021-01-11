package models

type Recipe struct {
	Id               string             `json:"id"`
	author           int                `json:"author"`
	contributors     map[int]string     `json:"contributors"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	IngredientGroups []IngredientGroups `json:"ingredientGroups"`
	Ingredients      []Ingredients      `json:"ingredients"`
	Qty              string             `json:"qty"`
	Unit             string             `json:"unit"`
	Steps            []Steps            `json:"steps"`
}
type IngredientGroups struct {
	GroupName string `json:"groupName"`
}
type Ingredients struct {
	Name string `json:"name"`
}
type Steps struct {
	ID      int               `json:"id"`
	Type    string            `json:"type"`
	Content map[string]string `json:"content"`
}
