package models

type Cocktail struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Image        string   `json:"image"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
	ABV          int      `json:"abv"`          
	ServingSize  string   `json:"servingSize"`  
}