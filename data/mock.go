package data

import "karl-s-bar-api/models"

var MockCocktails = []models.Cocktail{
	{
		ID:          "1",
		Name:        "Mojito",
		Category:    "Cocktail",
		Description: "Refreshing Cuban cocktail with rum, mint, and lime",
		Image:       "https://images.pexels.com/photos/3407750/pexels-photo-3407750.jpeg?w=800&h=600&fit=crop",
		Ingredients: []string{
			"60 ml White Rum",
			"22 ml Fresh Lime Juice",
			"15 ml Simple Syrup",
			"8-10 Fresh Mint Leaves",
			"Club Soda",
			"Ice",
			"Lime Slice",
		},
		Instructions: "1. Add mint leaves and simple syrup to a tall glass. 2. Gently muddle the mint to release oils, being careful not to crush leaves. 3. Add white rum and lime juice. 4. Fill the glass with ice and top with club soda. 5. Stir gently and garnish with a sprig of mint and a lime slice.",
		ABV:          12,
		ServingSize:  "355 ml",
	},
	{
		ID:          "2",
		Name:        "Margarita",
		Category:    "Cocktail",
		Description: "Classic cocktail with tequila, triple sec, and lime juice",
		Image:       "https://images.pexels.com/photos/2270690/pexels-photo-2270690.jpeg?w=800&h=600&fit=crop",
		Ingredients: []string{
			"60 ml Tequila",
			"30 ml Cointreau/Triple Sec",
			"30 ml Fresh Lime Juice",
			"15 ml Agave Nectar",
			"Ice",
			"Lime wheel for garnish",
			"Salt for rimmer",
		},
		Instructions: "1. Rim a chilled cocktail glass with salt. 2. Fill glass with ice. 3. Pour tequila, triple sec, and lime juice into a shaker with ice. 4. Shake well for about 10-15 seconds. 5. Strain into prepared glass. 6. Garnish with lime wheel.",
		ABV:          22,
		ServingSize:  "135 ml",
	},
	{
		ID:          "3",
		Name:        "Old Fashioned",
		Category:    "Whiskey",
		Description: "Timeless classic with bourbon, bitters, and a twist",
		Image:       "https://images.pexels.com/photos/3962282/pexels-photo-3962282.jpeg?w=800&h=600&fit=crop",
		Ingredients: []string{
			"60 ml Bourbon Whiskey",
			"1 Sugar Cube",
			"2 Dashes Angostura Bitters",
			"Few drops Water",
			"Ice",
			"Orange twist",
		},
		Instructions: "1. Place sugar cube in an old fashioned glass. 2. Add bitters and a few drops of water. 3. Muddle gently to dissolve sugar. 4. Add a large ice cube. 5. Pour bourbon over ice. 6. Stir gently. 7. Express orange oils over the drink and add as garnish.",
		ABV:          40,
		ServingSize:  "60 ml",
	},
	{
		ID:          "4",
		Name:        "Piña Colada",
		Category:    "Cocktail",
		Description: "Tropical blend of rum, coconut cream, and pineapple",
		Image:       "https://images.pexels.com/photos/3407753/pexels-photo-3407753.jpeg?w=800&h=600&fit=crop",
		Ingredients: []string{
			"60 ml White Rum",
			"90 ml Fresh Pineapple Juice",
			"45 ml Coconut Cream",
			"Ice",
			"Pineapple wedge",
			"Maraschino cherry",
		},
		Instructions: "1. Fill a cocktail shaker with ice. 2. Add rum, pineapple juice, and coconut cream. 3. Shake vigorously for 10-15 seconds. 4. Strain into a chilled glass filled with ice. 5. Garnish with pineapple wedge and cherry.",
		ABV:          11,
		ServingSize:  "195 ml",
	},
	{
		ID:          "5",
		Name:        "Cosmopolitan",
		Category:    "Cocktail",
		Description: "Vodka-based cocktail with cranberry juice and lime",
		Image:       "https://images.pexels.com/photos/3962314/pexels-photo-3962314.jpeg?w=800&h=600&fit=crop",
		Ingredients: []string{
			"45 ml Vodka",
			"30 ml Cointreau",
			"15 ml Fresh Lime Juice",
			"45 ml Cranberry Juice",
			"Ice",
			"Lime twist",
		},
		Instructions: "1. Fill a cocktail shaker with ice. 2. Add vodka, Cointreau, lime juice, and cranberry juice. 3. Shake well for 10-15 seconds. 4. Strain into a chilled martini glass. 5. Garnish with a lime twist.",
		ABV:          18,
		ServingSize:  "135 ml",
	},
	{
		ID:          "6",
		Name:        "Espresso Martini",
		Category:    "Cocktail",
		Description: "Energizing blend of vodka, coffee liqueur, and espresso",
		Image:       "https://images.pexels.com/photos/3962320/pexels-photo-3962320.jpeg?w=800&h=600&fit=crop",
		Ingredients: []string{
			"60 ml Vodka",
			"15 ml Kahlúa / Coffee Liqueur",
			"30 ml Fresh Espresso",
			"Ice",
			"Coffee beans for garnish",
		},
		Instructions: "1. Brew a fresh shot of espresso and let it cool slightly. 2. Fill a cocktail shaker with ice. 3. Add vodka and Kahlúa. 4. Add fresh espresso to shaker. 5. Shake vigorously for 10-15 seconds. 6. Strain into a chilled martini glass. 7. Garnish with 3 coffee beans on top.",
		ABV:          20,
		ServingSize:  "105 ml",
	},
}

type CocktailGetterImpl struct {}

func (c *CocktailGetterImpl)GetCocktails() ([]models.Cocktail, error){
	return MockCocktails, nil
}
	
func (c *CocktailGetterImpl) GetCocktailByID(id string) (models.Cocktail, error) {
	for _, cocktail := range MockCocktails {
		if cocktail.ID == id {
			return cocktail, nil
		}
	}
	return models.Cocktail{}, nil
}