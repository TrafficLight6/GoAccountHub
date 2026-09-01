package jsonOperate

import (
	"encoding/json"
	"fmt"
)

func Unmarshal[T any](jsonString string, result T, consoleAddend string) {
	if jsonString == "" {
		fmt.Println("⚠️ JSON String is Empty" + consoleAddend)
		return
	}
	if err := json.Unmarshal([]byte(jsonString), result); err != nil {
		fmt.Println("⚠️ JSON Unmarshal Error: " + err.Error() + consoleAddend)
	}
}

func Marshal[T any](input T) string {
	bytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("⚠️ JSON Marshal Error: " + err.Error())
		return ""
	}
	return string(bytes)
}
