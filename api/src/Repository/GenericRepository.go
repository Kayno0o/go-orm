package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
	utils "go-api-test.kayn.ooo/src/Utils"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	DB                *bun.DB
	Ctx               = context.Background()
	GenericRepository = GenericRepositoryStruct{}
)

type GenericRepositoryInterface interface {
	FindOneById(entity interface{}, id uint) error
	FindOneBy(entity interface{}, params map[string]interface{}) error
	FindAll(entities interface{}) error
	FindAllBy(entities interface{}, params map[string]interface{}) error
	CountAll(entity interface{}) (int, error)
	Create(entity interface{}) (sql.Result, error)
	Update(entity interface{}) (sql.Result, error)
}

type GenericRepositoryStruct struct {
	GenericRepositoryInterface
}

func (r *GenericRepositoryStruct) Init(entities []interface{}) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("DB_URL environment variable is required")
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

func (r *GenericRepositoryStruct) FindOneById(entity interface{}, id uint) error {
	return DB.NewSelect().Model(entity).Where("id = ?", id).Scan(Ctx)
}

func (r *GenericRepositoryStruct) applyParams(model *bun.SelectQuery, params map[string]interface{}) *bun.SelectQuery {
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

func (r *GenericRepositoryStruct) FindOneBy(entity interface{}, params map[string]interface{}) error {
	model := DB.NewSelect().Model(entity)
	params["limit"] = 1
	params["offset"] = 0
	model = r.applyParams(model, params)
	return model.Scan(Ctx)
}

func (r *GenericRepositoryStruct) FindAll(entities interface{}) error {
	return DB.NewSelect().Model(entities).OrderExpr("id ASC").Scan(Ctx)
}

func (r *GenericRepositoryStruct) FindAllBy(entities interface{}, params map[string]interface{}) error {
	model := DB.NewSelect().Model(entities)
	model = r.applyParams(model, params)
	return model.OrderExpr("id ASC").Scan(Ctx)
}

func (r *GenericRepositoryStruct) CountAll(entity interface{}) (int, error) {
	return DB.NewSelect().Model(entity).Count(Ctx)
}

func (r *GenericRepositoryStruct) Create(entity interface{}) (sql.Result, error) {
	return DB.NewInsert().Model(entity).Exec(Ctx)
}

func (r *GenericRepositoryStruct) Update(entity interface{}) (sql.Result, error) {
	if timestampableEntity, ok := entity.(trait.Timestampable); ok {
		timestampableEntity.UpdatedAt = time.Now()
	}
	model := DB.NewUpdate().Model(entity)
	if identifiableEntity, ok := entity.(trait.IdentifierInterface); ok {
		model.Where("id = ?", identifiableEntity.GetId())
	}

	return model.Exec(Ctx)
}

func FindEntityByRouteParam(c *fiber.Ctx, param string, entity interface{}) error {
	id, err := strconv.ParseUint(c.Params(param), 10, 0)
	if err != nil {
		return utils.HTTP404Error(c)
	}

	err = GenericRepository.FindOneById(entity, uint(id))
	if entity == nil || err != nil {
		return utils.HTTP404Error(c)
	}

	return nil
}
