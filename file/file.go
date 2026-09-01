package file

import (
	"fmt"
	"os"
)

func Read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("⚠️ File Read Error: " + err.Error() + " on " + path)
		return ""
	}
	return string(data)
}

func Write(path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Println("⚠️ File Write Error: " + err.Error() + " on " + path)
	}
}
