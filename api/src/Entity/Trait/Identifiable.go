package trait

type IdentifiableTraitI interface {
	GetId() uint
}

type IdentifiableTrait struct {
	IdentifiableTraitI `bun:"-"`

	ID uint `bun:",pk,autoincrement" json:"id"`
}

func (i IdentifiableTrait) GetId() uint {
	return i.ID
}
