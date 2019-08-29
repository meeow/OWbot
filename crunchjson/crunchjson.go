package crunchjson

import (
	"encoding/json"
)

// JSONtoMap takes JSON blob and returns corresponding map
func JSONtoMap(jsonString string) map[string]interface{} {
	jsonBytes := []byte(jsonString)
	var f interface{}
	json.Unmarshal(jsonBytes, &f)
	mapping := f.(map[string]interface{})

	return mapping
}
