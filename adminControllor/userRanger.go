package adminControllor

import (
	"net/http"
	"strings"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRangeRequestBody struct {
	//It the Length is -1, It Means Get All User from BeginTableId to The End
	BeginTableId    int                 `json:"begin_table_id"`
	Length          int                 `json:"length"`
	SearchCondition UserSearchCondition `json:"search_condition"`
}

type UserSearchCondition struct {
	//If The Field is Empty, It Means Ignore This Condition
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	UUHash       string `json:"uu_hash"`
}

func UserRange(c *gin.Context) {
	var body UserRangeRequestBody
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
	//Build Search Query (only non-empty fields, combined with AND)
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
	query := db.Model(&sqlTable.User{})
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	//Search
	var users []sqlTable.User
	if isReturnAll {
		query.Find(&users)
	} else {
		query.Limit(body.Length).Offset((body.BeginTableId - 1) * body.Length).Find(&users)
	}
	//Return Result
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success", "data": users})
	return
}
