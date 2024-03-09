package fixture

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	entity "go-api-test.kayn.ooo/src/Entity"
	repository "go-api-test.kayn.ooo/src/Repository"
	security "go-api-test.kayn.ooo/src/Security"
)

func GetFirstNames() []string {
	file, err := os.ReadFile("./firstnames.txt")
	if err != nil {
		panic(err)
	}

	return strings.Split(string(file), "\n")
}

func RandomFirstName(firstNames []string) string {
	return firstNames[rand.Intn(len(firstNames))]
}

func GenerateUsersFromJson() []entity.User {
	var users []entity.User

	data, err := os.ReadFile("./src/Fixture/users.json")
	if err != nil {
		fmt.Println("Error reading users.json:", err)
		return users
	}

	var usersData []struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	if err := json.Unmarshal(data, &usersData); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return users
	}

	for _, userData := range usersData {
		role := "ROLE_USER"
		if userData.IsAdmin {
			role = "ROLE_ADMIN"
		}

		user := entity.User{
			Username: userData.Username,
			Email:    userData.Email,
			Password: security.HashPassword(userData.Password),
			Roles:    []string{role},
		}

		users = append(users, user)
	}

	_, err = repository.DB.NewInsert().Model(&users).Exec(repository.Ctx)
	if err != nil {
		panic(err)
	}

	return users
}

func GenerateUsers(nb int, isAdmin bool) []entity.User {
	password := security.HashPassword("password")

	var users []entity.User
	firstNames := GetFirstNames()
	for i := 0; i < nb; i++ {
		firstName := RandomFirstName(firstNames)

		role := "ROLE_USER"
		if isAdmin {
			role = "ROLE_ADMIN"
		}

		user := &entity.User{
			Username: firstName,
			Email:    firstName + "@gmail.com",
			Password: password,
			Roles:    []string{role},
		}

		users = append(users, *user)
	}

	_, err := repository.DB.NewInsert().Model(&users).Exec(repository.Ctx)
	if err != nil {
		panic(err)
	}

	return users
}
