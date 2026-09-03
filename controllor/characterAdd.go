package controllor

import "github.com/gin-gonic/gin"

type CharacterAddRequestBody struct {
	CharacterName string `json:"character_name"`
	Password      string `json:"password"`
	UserUUHash    string `json:"user_uu_hash"`
	MetaData      string `json:"meta_data"`
}

func CharacterAdd(c *gin.Context) {
}
