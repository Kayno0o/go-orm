package routerstruct

import entity "go-api-test.kayn.ooo/src/Entity"

type Route struct {
	Title  string
	Params map[string]string
	Query  map[string]string
	User   *entity.User
}
