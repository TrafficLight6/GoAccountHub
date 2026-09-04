package adminControllor

import (
	"net/http"
	"strings"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminRangeRequestBody struct {
	//It the Length is -1, It Means Get All Admin from BeginTableId to The End
	BeginTableId    int             `json:"begin_table_id"`
	Length          int             `json:"length"`
	SearchCondition SearchCondition `json:"search_condition"`
}

type SearchCondition struct {
	//If The Field is Empty, It Means Ignore This Condition
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	UUHash       string `json:"uu_hash"`
}

func AdminRange(c *gin.Context) {
	var body AdminRangeRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		c.Abort()
		return
	}
	if body.BeginTableId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "BeginTableId Begin From 1"})
		c.Abort()
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
	//Build Search Query (only non-empty fields, combined with OR)
	condition := body.SearchCondition
	var conditions []string
	var args []interface{}
	if condition.Username != "" {
		conditions = append(conditions, "username LIKE ?")
		args = append(args, "%"+condition.Username+"%")
	}
	if condition.PasswordHash != "" {
		conditions = append(conditions, "password_hash LIKE ?")
		args = append(args, "%"+condition.PasswordHash+"%")
	}
	if condition.UUHash != "" {
		conditions = append(conditions, "uu_hash LIKE ?")
		args = append(args, "%"+condition.UUHash+"%")
	}
	query := db.Model(&sqlTable.Admin{})
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}
	//Search
	var admins []sqlTable.Admin
	if isReturnAll {
		query.Find(&admins)
	} else {
		query.Limit(body.Length).Offset((body.BeginTableId - 1) * body.Length).Find(&admins)
	}
	//Return Result
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success(not include root admin)", "data": admins})
}
