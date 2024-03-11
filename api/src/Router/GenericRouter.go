package router

import (
	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
	middleware "go-api-test.kayn.ooo/src/Middleware"
	repository "go-api-test.kayn.ooo/src/Repository"
	utils "go-api-test.kayn.ooo/src/Utils"
	"strconv"
)

var FiberApp *fiber.App

type GenericRouterI interface {
	RegisterRoutes(fiber.Router)
}

type Params struct {
	VerifyOwner bool
	Pagination  bool
}

type CrudParams struct {
	PublicGet     bool
	PublicContext interface{}
}

func GetQuery[C trait.IdentifiableTraitI](c *fiber.Ctx, params Params) map[string]interface{} {
	query := make(map[string]interface{})

	if params.VerifyOwner {
		uid := utils.GetUserId(c)
		if any(uid) == nil {
			return nil
		}

		var e C
		if _, ok := any(e).(entity.OwnerableTraitI); ok {
			query["owner_id"] = uid
		}
	}

	if params.Pagination {
		query = utils.MergeMaps(query, GetPagination(c))
	}

	return query
}

func GetPagination(c *fiber.Ctx) map[string]interface{} {
	query := make(map[string]interface{})

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}
	query["limit"] = limit

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	query["offset"] = (page - 1) * limit

	return query
}

func GetOne[E trait.IdentifiableTraitI, C interface{}](
	params Params,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		query := GetQuery[E](c, params)
		if query == nil {
			return utils.HTTP404Error(c)
		}

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return utils.HTTP400Error(c, err.Error())
		}
		query["id"] = id

		var e E
		e, err = repository.FindOneBy[E](query)
		if err != nil || any(e) == nil || e.GetId() == 0 {
			return utils.HTTP404Error(c)
		}

		return utils.JsonContext[C](c, e)
	}
}

func GetAll[E trait.IdentifiableTraitI, C interface{}](
	params Params,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		query := GetQuery[E](c, params)
		if query == nil {
			return utils.HTTP404Error(c)
		}

		entities, err := repository.FindAllBy[E](query)
		if err != nil {
			return utils.HTTP500Error(c, err.Error())
		}

		return utils.JsonContext[[]C](c, entities)
	}
}

func CountAll[E trait.IdentifiableTraitI](
	params Params,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		query := GetQuery[E](c, params)
		if query == nil {
			return utils.HTTP404Error(c)
		}

		count, err := repository.CountAllBy[E](query)
		if err != nil {
			return utils.HTTP500Error(c, err.Error())
		}
		return c.JSON(count)
	}
}

func Put[E trait.IdentifiableTraitI, Input interface{}, Output interface{}](
	params Params,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		query := GetQuery[E](c, params)
		if query == nil {
			return utils.HTTP404Error(c)
		}
		if c.Params("id") == "" {
			return utils.HTTP404Error(c)
		}

		e, err := repository.FindEntityByRouteParam[E](c, "id")
		if err != nil {
			return nil
		}

		var input Input
		if err := c.BodyParser(&input); err != nil {
			return utils.HTTP400Error(c, err.Error())
		}

		if params.VerifyOwner {
			if !utils.IsOwner(c, e) {
				return utils.HTTP401Error(c)
			}
		}

		err = utils.ApplyEntity(&e, input)
		if err != nil {
			return utils.HTTP400Error(c, err.Error())
		}

		_, err = repository.Update(&e, e.GetId())
		if err != nil {
			return utils.HTTP400Error(c, err.Error())
		}

		return utils.JsonContext[Output](c, e)
	}
}

func Post[E trait.IdentifiableTraitI, C interface{}, OC interface{}](
	create func(*fiber.Ctx, C) (E, error),
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var input C
		if err := c.BodyParser(&input); err != nil {
			return utils.HTTP400Error(c, err.Error())
		}

		uid := utils.GetUserId(c)
		if any(uid) == nil {
			return utils.HTTP401Error(c)
		}

		e, err := create(c, input)
		if err != nil {
			return utils.HTTP400Error(c, err.Error())
		}

		_, err = repository.Create(&e)
		if err != nil {
			return utils.HTTP400Error(c, err.Error())
		}

		return utils.JsonContext[OC](c, e)
	}
}

func Delete[E trait.IdentifiableTraitI](
	params Params,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		query := GetQuery[E](c, params)
		if query == nil {
			return utils.HTTP404Error(c)
		}
		if c.Params("id") == "" {
			return utils.HTTP404Error(c)
		}

		e, err := repository.FindEntityByRouteParam[E](c, "id")
		if err != nil {
			return nil
		}

		if params.VerifyOwner {
			if !utils.IsOwner(c, e) {
				return utils.HTTP401Error(c)
			}
		}

		_, err = repository.Delete(&e)
		if err != nil {
			return utils.HTTP400Error(c)
		}
		return nil
	}
}

func RegisterCrud[E trait.IdentifiableTraitI, C interface{}, UC interface{}](
	r fiber.Router,
	name string,
	params CrudParams,
	createHandler func(*fiber.Ctx, UC) (E, error),
) {
	adminParams := Params{VerifyOwner: false, Pagination: true}
	userParams := Params{VerifyOwner: true, Pagination: true}
	entityParams := Params{VerifyOwner: !params.PublicGet, Pagination: true}

	r.Group(
		"/admin/"+name,
		middleware.IsGranted([]string{"ROLE_ADMIN"}),
	).Post(
		"/",
		Post[E, UC, C](createHandler),
	).Put(
		"/:id",
		Put[E, E, C](adminParams),
	).Delete(
		"/:id",
		Delete[E](adminParams),
	).Get(
		"/:id",
		GetOne[E, C](adminParams),
	).Get(
		"/",
		GetAll[E, C](adminParams),
	).Get(
		"/count",
		CountAll[E](adminParams),
	)

	r.Group(
		name,
		middleware.IsGranted([]string{"ROLE_USER"}),
	).Post(
		"/",
		Post[E, UC, C](createHandler),
	).Put(
		"/:id",
		Put[E, UC, C](userParams),
	).Delete(
		"/:id",
		Delete[E](userParams),
	)

	r.Group(name).Get(
		"/:id",
		GetOne[E, C](entityParams),
	).Get(
		"/",
		GetAll[E, C](entityParams),
	).Get(
		"/count",
		CountAll[E](entityParams),
	)
}
