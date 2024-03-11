package trait

import "time"

type TimestampableTraitI interface {
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}

type TimestampableTrait struct {
	TimestampableTraitI `bun:"-"`

	CreatedAt time.Time `bun:",nullzero,default:now()"`
	UpdatedAt time.Time `bun:",nullzero,default:now()"`
}

func (t *TimestampableTrait) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *TimestampableTrait) SetCreatedAt(at time.Time) {
	t.CreatedAt = at
}

func (t *TimestampableTrait) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}

func (t *TimestampableTrait) SetUpdatedAt(at time.Time) {
	t.UpdatedAt = at
}
