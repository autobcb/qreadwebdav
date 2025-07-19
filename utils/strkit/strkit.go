package strkit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ToString(obj interface{}) string {
	if obj == nil {
		return ""
	}
	switch obj.(type) {
	case string:
		return obj.(string)
	case int64:
		return strconv.FormatInt(obj.(int64), 10)
	case int:
		return ToString(int64(obj.(int)))
	case int8:
		return ToString(int64(obj.(int8)))
	case int16:
		return ToString(int64(obj.(int16)))
	case int32:
		return ToString(int64(obj.(int32)))
	case float64:
		return strconv.FormatFloat(obj.(float64), 'f', -1, 64)
	case float32:
		return ToString(int64(obj.(float32))) //strconv.FormatFloat(float64(obj.(float32)), 'f', -1, 64)
	case bool:
		if obj.(bool) {
			return "true"
		} else {
			return "false"
		}
	case []byte:
		return string(obj.([]byte))
	case byte:
		return string(obj.(byte))
	default:
		return ToJson(obj)
	}

}

func ToJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

func GetStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		return val.(string)
	}
	return ""
}

func SanitizeFilename(filename string) string {
	invalidChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "_")
	}

	const maxLength = 50
	if len(filename) > maxLength {
		prefix := filename[:20]
		suffix := filename[len(filename)-30:]
		filename = prefix + "..." + suffix
	}
	return filename
}

func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}
	return result, nil
}
