package router

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
	middleware "go-api-test.kayn.ooo/src/Middleware"
	repository "go-api-test.kayn.ooo/src/Repository"
	utils "go-api-test.kayn.ooo/src/Utils"
)

type GenericRouterInterface interface {
	RegisterRoutes(fiber.Router)
}

func Init(routers []GenericRouterInterface) {
	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	api := fiberApp.Group("/", middleware.Authenticate)
	for i := range routers {
		routers[i].RegisterRoutes(api)
	}

	log.Fatal(fiberApp.Listen(":3000"))
}

func queryToParams(c *fiber.Ctx) map[string]interface{} {
	params := map[string]interface{}{}
	for key, value := range c.Queries() {
		params[key] = value
	}
	return params
}

func FindOne(rep repository.GenericRepositoryInterface, entity trait.IdentifierInterface, context interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return utils.HTTP400Error(c)
		}
		_ = rep.FindOneById(entity, uint(id))
		if entity.GetId() == 0 {
			return utils.HTTP404Error(c)
		}
		return utils.JsonContext(c, entity, context)
	}
}

func FindAll(rep repository.GenericRepositoryInterface, entities interface{}, context interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		params := queryToParams(c)
		if params["offset"] == nil {
			params["offset"] = 1
		}
		if params["limit"] == nil {
			params["limit"] = 10
		}

		err := rep.FindAllBy(entities, params)
		if err != nil {
			return utils.HTTP500Error(c)
		}
		return utils.JsonContext(c, entities, context)
	}
}

func CountAll(rep repository.GenericRepositoryInterface, entity trait.IdentifierInterface) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		count, err := rep.CountAll(entity)
		if err != nil {
			return utils.HTTP500Error(c)
		}
		return c.JSON(count)
	}
}
