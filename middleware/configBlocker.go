package middleware

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/gin-gonic/gin"
)

func ConfigBlocker(field string) gin.HandlerFunc {
	return func(c *gin.Context) {
		config := c.Value("config").(config.Config)
		switchConfig := config.SwitchConfig
		items := parseSwitchConfig(switchConfig)
		if !items[field] {
			c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "error": ErrorMessage(field)})
			c.Abort()
			return
		}
		c.Next()
	}
}

// parseSwitchConfig parses SwitchConfig to map of bool switches (e.g. allow_multi_character)
func parseSwitchConfig(switchConfig config.SwitchConfig) map[string]bool {
	items := make(map[string]bool)
	typeOf := reflect.TypeOf(switchConfig)
	valueOf := reflect.ValueOf(switchConfig)
	for i := 0; i < typeOf.NumField(); i++ {
		structField := typeOf.Field(i)
		name := strings.Split(structField.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			name = structField.Name
		}
		items[name] = valueOf.Field(i).Bool()
	}
	return items
}

func ErrorMessage(field string) string {
	if field == "allow_multi_character" {
		return "Server Not Allow Multi Character"
	}
	return "Switch Disabled: " + field
}
