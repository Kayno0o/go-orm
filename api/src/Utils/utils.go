package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/gofiber/fiber/v2"
	entity "go-api-test.kayn.ooo/src/Entity"
	trait "go-api-test.kayn.ooo/src/Entity/Trait"
)

func Includes(array []string, search string) bool {
	for _, element := range array {
		if element == search {
			return true
		}
	}

	return false
}

func ApplyContextInto[I any, O any](input I, output *O) error {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonInput, &output); err != nil {
		return err
	}

	return nil
}

// ApplyContext takes and input entity and output json format
func ApplyContext[C interface{}](input interface{}) (C, error) {
	var c C
	err := ApplyContextInto(input, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func JsonContext[C interface{}](c *fiber.Ctx, input interface{}) error {
	context, err := ApplyContext[C](input)
	if err != nil {
		return HTTP400Error(c, err.Error())
	}
	return c.JSON(context)
}

func MapToStruct[T any](m map[string]any) (T, error) {
	var result T
	mJSON, err := json.Marshal(m)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(mJSON, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ApplyEntity[T trait.IdentifiableTraitI](output *T, input interface{}) error {
	var m map[string]any

	err := ApplyContextInto(output, &m)
	if err != nil {
		return err
	}

	err = ApplyContextInto(input, &m)
	if err != nil {
		return err
	}

	jsonOutput, err := json.Marshal(m)

	if err := json.Unmarshal(jsonOutput, &output); err != nil {
		return err
	}

	return err
}

func GetUserId(c *fiber.Ctx) uint {
	return GetUser(c).ID
}

func IsOwner[E trait.IdentifiableTraitI](c *fiber.Ctx, e E) bool {
	if u, ok := any(e).(entity.User); ok {
		return u.ID == GetUserId(c)
	}

	if oe, ok := any(e).(entity.OwnerableTraitI); ok {
		return oe.GetOwnerId() == GetUserId(c)
	} else {
		return false
	}
}

func MergeMaps[U comparable, T any](maps ...map[U]T) map[U]T {
	mergedMap := make(map[U]T)

	for _, val := range maps {
		for key, value := range val {
			mergedMap[key] = value
		}
	}

	return mergedMap
}

func GetUser[T *entity.User](c *fiber.Ctx) T {
	return c.Locals("user").(T)
}

func Stringify(r any) string {
	str, err := json.Marshal(&r)
	if err != nil {
		return ""
	}
	return string(str)
}

func GetHexColor(hex string) (string, error) {
	re := regexp.MustCompile(`(?m)^(?:\#|0x|)([a-fA-F0-9]{6})$`)
	matches := re.FindStringSubmatch(hex)
	if len(matches) > 1 {
		return "#" + matches[1], nil
	}
	return "", errors.New("No color matches")
}

func MapToArray[T any, U comparable](m map[U]T) []T {
	items := make([]T, 0)
	for _, value := range m {
		items = append(items, value)
	}
	return items
}

func RandomString(length int) (string, error) {
	// Create a byte slice to hold the random bytes
	randomBytes := make([]byte, length)

	// Read random bytes from crypto/rand
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode random bytes to a base64 string
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	return randomString, nil
}
