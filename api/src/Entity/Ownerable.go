package entity

type OwnerableTraitI interface {
	GetOwnerId() uint
}

type OwnerableTrait struct {
	OwnerableTraitI `bun:"-"`

	OwnerId uint `bun:",notnull"`
}

func (o OwnerableTrait) GetOwnerId() uint {
	return o.OwnerId
}
