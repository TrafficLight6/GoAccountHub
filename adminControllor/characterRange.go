package adminControllor

import (
	"net/http"
	"strings"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CharacterRangeRequestBody struct {
	//It the Length is -1, It Means Get All User from BeginTableId to The End
	BeginTableId    int                      `json:"begin_table_id"`
	Length          int                      `json:"length"`
	SearchCondition CharacterSearchCondition `json:"search_condition"`
}

type CharacterSearchCondition struct {
	//If The Field is Empty, It Means Ignore This Condition
	CharacterName string `json:"character_name"`
	PasswordHash  string `json:"password_hash"`
	UserUUHash    string `json:"user_uu_hash"`
	UUHash        string `json:"uu_hash"`
}

func CharacterRange(c *gin.Context) {
	var body CharacterRangeRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
		return
	}
	//Check Length is Valid
	if body.Length < -1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Length Must Be Greater Than -1"})
		c.Abort()
		return
	}
	var isReturnAll bool
	if body.Length == -1 {
		isReturnAll = true
	} else {
		isReturnAll = false
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Build Search Query (only non-empty fields, combined with AND)
	condition := body.SearchCondition
	var conditions []string
	var args []interface{}
	if condition.CharacterName != "" {
		conditions = append(conditions, "character_name LIKE ?")
		args = append(args, "%"+condition.CharacterName+"%")
	}
	if condition.PasswordHash != "" {
		conditions = append(conditions, "password_hash LIKE ?")
		args = append(args, "%"+condition.PasswordHash+"%")
	}
	if condition.UserUUHash != "" {
		conditions = append(conditions, "user_uu_hash LIKE ?")
		args = append(args, "%"+condition.UserUUHash+"%")
	}
	if condition.UUHash != "" {
		conditions = append(conditions, "uu_hash LIKE ?")
		args = append(args, "%"+condition.UUHash+"%")
	}
	query := db.Model(&sqlTable.Character{})
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	//Search
	var characters []sqlTable.Character
	if isReturnAll {
		query.Find(&characters)
	} else {
		query.Limit(body.Length).Offset((body.BeginTableId - 1) * body.Length).Find(&characters)
	}
	//Return Result
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success", "data": characters})
	return
}
