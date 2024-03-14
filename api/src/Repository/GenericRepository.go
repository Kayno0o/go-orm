package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
	utils "go-api-test.kayn.ooo/src/Utils"
)

var (
	DB  *bun.DB
	Ctx = context.Background()
)

type GenericRepositoryStruct[E trait.IdentifiableTraitI] struct {
}

func Init(entities []interface{}) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("DB_URL environment variable is required")
		os.Exit(1)
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbURL)))

	DB = bun.NewDB(sqldb, pgdialect.New())

	file, err := os.Create("bundebug.log")
	if err != nil {
		panic(err)
	}

	DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.WithWriter(file),
	))

	if err := DB.Ping(); err != nil {
		panic(err)
	}

	for i := range entities {
		DB.RegisterModel(entities[i])
		_, err := DB.NewCreateTable().Model(entities[i]).IfNotExists().Exec(Ctx)
		if err != nil {
			panic(err)
		}
	}
}

func applyParams(model *bun.SelectQuery, params map[string]interface{}) *bun.SelectQuery {
	for key, value := range params {
		if key == "limit" || key == "offset" {
			limit, boolErr := value.(string)
			if !boolErr {
				continue
			}

			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				continue
			}

			if key == "offset" {
				model.Offset(limitInt)
			} else if key == "limit" {
				model.Limit(limitInt)
			}
			continue
		}

		regex := regexp.MustCompile(`^[a-z_]+$`)
		if !regex.MatchString(key) {
			continue
		}
		model.Where(key+" = ?", value)
	}
	return model
}

func FindEntityByRouteParam[E trait.IdentifiableTraitI](c *fiber.Ctx, param string) (E, error) {
	var e E
	id, err := strconv.ParseUint(c.Params(param), 10, 0)
	if err != nil {
		return e, utils.HTTP404Error(c, err.Error())
	}

	e, err = FindOneById[E](uint(id))
	if err != nil {
		return e, utils.HTTP404Error(c, err.Error())
	}

	if any(e) == nil {
		return e, utils.HTTP404Error(c)
	}

	return e, nil
}

func FindOneById[E trait.IdentifiableTraitI](id uint) (E, error) {
	var e E
	err := DB.NewSelect().Model(&e).Where("id = ?", id).Scan(Ctx)
	return e, err
}

func FindOneBy[E trait.IdentifiableTraitI](params map[string]interface{}) (E, error) {
	var e E
	model := DB.NewSelect().Model(&e)
	params["limit"] = 1
	params["offset"] = 0
	model = applyParams(model, params)
	err := model.Scan(Ctx)
	return e, err
}

func FindAll[E trait.IdentifiableTraitI]() ([]E, error) {
	var entities []E
	err := DB.NewSelect().Model(&entities).OrderExpr("id ASC").Scan(Ctx)
	return entities, err
}

func FindAllBy[E trait.IdentifiableTraitI](params map[string]interface{}) ([]E, error) {
	var entities []E
	model := DB.NewSelect().Model(&entities)
	model = applyParams(model, params)
	err := model.OrderExpr("id ASC").Scan(Ctx)
	return entities, err
}

func CountAll[E trait.IdentifiableTraitI]() (int, error) {
	var e E
	return DB.NewSelect().Model(&e).Count(Ctx)
}

func CountAllBy[E trait.IdentifiableTraitI](params map[string]interface{}) (int, error) {
	var e E
	model := DB.NewSelect().Model(&e)
	model = applyParams(model, params)
	return model.Count(Ctx)
}

func Create[E trait.IdentifiableTraitI](e *E) (sql.Result, error) {
	model := DB.NewInsert().Model(e)
	if timestampableEntity, ok := any(e).(trait.TimestampableTraitI); ok {
		timestampableEntity.SetCreatedAt(time.Now())
	}

	return model.Exec(Ctx)
}

func Update[E any](e *E, id uint) (sql.Result, error) {
	if timestampableEntity, ok := any(e).(trait.TimestampableTraitI); ok {
		timestampableEntity.SetUpdatedAt(time.Now())
	}
	model := DB.NewUpdate().Model(e)
	model.Where("id = ?", id)

	return model.Exec(Ctx)
}

func Delete[E trait.IdentifiableTraitI](e *E) (sql.Result, error) {
	model := DB.NewDelete().Model(e)
	model.Where("id = ?", (*e).GetId())

	return model.Exec(Ctx)
}
