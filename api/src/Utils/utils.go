package utils

import "encoding/json"

func Includes(array []string, search string) bool {
	for _, element := range array {
		if element == search {
			return true
		}
	}

	return false
}

// ApplyContext takes and input entity and output json format
func ApplyContext(input interface{}, context interface{}) {
	jsonInput, _ := json.Marshal(input)
	_ = json.Unmarshal(jsonInput, context)
}
