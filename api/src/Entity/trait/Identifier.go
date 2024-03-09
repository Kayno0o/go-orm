package trait

type IdentifierInterface interface {
	GetId() uint
}

type Identifier struct {
	IdentifierInterface `bun:"-" json:"-"`

	ID uint `bun:",pk,autoincrement" json:"id"`
}

func (i *Identifier) GetId() uint {
	return i.ID
}
