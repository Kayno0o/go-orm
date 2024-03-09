package router

import (
	middleware "go-api-test.kayn.ooo/src/Middleware"
	utils "go-api-test.kayn.ooo/src/Utils"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	fixture "go-api-test.kayn.ooo/src/Fixture"
	repository "go-api-test.kayn.ooo/src/Repository"
	security "go-api-test.kayn.ooo/src/Security"
)

type UserRouter struct {
	GenericRouterInterface
}

func (ur *UserRouter) RegisterRoutes(r fiber.Router) {
	api := r.Group("/api")

	api.Post(
		"/user/login",
		ur.Login,
	).Post(
		"/user/register",
		ur.Register,
	)

	// ADMIN
	api.Get(
		"/users/fixture/:amount",
		middleware.IsGranted([]string{"ROLE_ADMIN"}),
		ur.Fixture,
	)

	// USER
	api.Get(
		"/user/me",
		middleware.IsGranted([]string{"ROLE_USER"}),
		ur.Me,
	)

	// PUBLIC
	api.Get(
		"/users",
		FindAll(
			repository.UserRepository,
			&[]entity.User{},
			&[]entity.UserContext{},
		),
	).Get(
		"/users/count",
		CountAll(
			repository.UserRepository,
			&entity.User{},
		),
	).Get(
		"/user/:id",
		FindOne(
			repository.UserRepository,
			&entity.User{},
			&entity.UserContext{},
		),
	)
}

func (ur *UserRouter) Login(c *fiber.Ctx) error {
	var login entity.Login
	if err := c.BodyParser(&login); err != nil {
		return utils.HTTP400Error(c)
	}

	user, err := security.Authenticate(&login)
	if err != nil {
		return utils.HTTP401Error(c)
	}

	token, err := security.GenerateToken(user)
	if err != nil {
		return utils.HTTP500Error(c)
	}

	return c.JSON(token)
}

func (ur *UserRouter) Register(c *fiber.Ctx) error {
	var form entity.Register
	if err := c.BodyParser(&form); err != nil {
		return utils.HTTP400Error(c)
	}

	var user entity.User
	user.Username = form.Username
	user.Email = form.Email

	password := security.HashPassword(form.Password)
	user.Password = password

	_, err := repository.UserRepository.Create(&user)
	if err != nil {
		return utils.HTTP500Error(c)
	}

	token, err := security.GenerateToken(&user)
	if err != nil {
		return utils.HTTP500Error(c)
	}

	// add token to session/cookies
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token.Token,
		Path:     "/",
		Expires:  token.ExpiresAt,
		HTTPOnly: true,
		Domain:   os.Getenv("DOMAIN"),
		Secure:   true,
	})

	return c.JSON(token)
}

func (ur *UserRouter) Fixture(c *fiber.Ctx) error {
	amount, err := strconv.Atoi(c.Params("amount"))
	if err != nil {
		return utils.HTTP400Error(c)
	}

	users := fixture.GenerateUsers(amount, false)

	return c.JSON(users)
}

func (ur *UserRouter) Me(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return utils.HTTP401Error(c)
	}

	return c.JSON(user)
}

func (ur *UserRouter) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Domain:   os.Getenv("DOMAIN"),
		Secure:   true,
	})

	return c.SendStatus(200)
}
