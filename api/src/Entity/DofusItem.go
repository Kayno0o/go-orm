package entity

import (
	"github.com/uptrace/bun"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
)

type DofusItem struct {
	bun.BaseModel `bun:"table:dofus_item,alias:ditem"`

	trait.IdentifiableTrait
	trait.TimestampableTrait

	Name         string               `bun:",notnull" json:"name"`
	Level        uint8                `bun:",notnull" json:"level"`
	PriceHistory []*DofusPriceHistory `bun:"rel:has-many,join:id=item_id"`
	Ingredients  []*DofusIngredient   `bun:"rel:has-many,join:id=recipe_id"`
}
