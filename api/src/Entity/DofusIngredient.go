package entity

import "github.com/uptrace/bun"

type DofusIngredient struct {
	bun.BaseModel `bun:"table:dofus_ingredient,alias:dingredient"`

	IngredientID uint       `bun:",pk"`
	RecipeID     uint       `bun:",pk"`
	Ingredient   *DofusItem `bun:"rel:belongs-to,join:ingredient_id=id"`
	Recipe       *DofusItem `bun:"rel:belongs-to,join:recipe_id=id"`
}
