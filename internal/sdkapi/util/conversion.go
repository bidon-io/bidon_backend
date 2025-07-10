package util

import "fmt"

// ConvertToStringMap converts map[string]any to map[string]string
func ConvertToStringMap(input map[string]any) map[string]string {
	if input == nil {
		return nil
	}

	result := make(map[string]string, len(input))
	for key, value := range input {
		if value == nil {
			result[key] = ""
		} else {
			result[key] = fmt.Sprintf("%v", value)
		}
	}
	return result
}
