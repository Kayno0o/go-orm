package entity

import (
	"time"

	"github.com/uptrace/bun"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
)

type DofusPriceHistory struct {
	bun.BaseModel `bun:"table:dofus_price_history,alias:dph"`

	trait.IdentifiableTrait
	OwnerableTrait

	Price  uint      `bun:",notnull" json:"price"`
	ItemID uint      `bun:",notnull"`
	Date   time.Time `bun:",nullzero,default:now()" json:"date"`
}
