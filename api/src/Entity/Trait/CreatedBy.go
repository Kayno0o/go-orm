package trait

import entity "go-api-test.kayn.ooo/src/Entity"

type CreatedBy struct {
	CreatedById uint         `bun:",notnull" json:"created_by_id"`
	CreatedBy   *entity.User `bun:"rel:belongs-to,join:created_by_id=id" json:"created_by"`
}
