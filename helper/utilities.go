package helper

import (
	"encoding/json"
	"regexp"
)

func StripAllButNumbers(str string) string {
	// var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	var nonAlphanumericRegex = regexp.MustCompile(`[^0-9]+`)
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func StructJson(data interface{}) string {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonString)
}
