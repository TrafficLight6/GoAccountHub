package adminControllor

import (
	"net/http"
	"strings"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KeyRangeRequestBody struct {
	//It the Length is -1, It Means Get All Key from BeginTableId to The End
	BeginTableId    int                `json:"begin_table_id"`
	Length          int                `json:"length"`
	SearchCondition KeySearchCondition `json:"search_condition"`
}

type KeySearchCondition struct {
	//If The Field is Empty, It Means Ignore This Condition
	KeyUser string `json:"key_user"`
}

func KeyRange(c *gin.Context) {
	var body KeyRangeRequestBody
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
	if condition.KeyUser != "" {
		conditions = append(conditions, "key_user LIKE ?")
		args = append(args, "%"+condition.KeyUser+"%")
	}
	query := db.Model(&sqlTable.ApplicationKey{})
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	//Must Be Ordered, Otherwise Limit/Offset Follows Physical Row Order And Pages Are Not Stable
	query = query.Order("id ASC")
	//Search
	var keys []sqlTable.ApplicationKey
	if isReturnAll {
		query.Find(&keys)
	} else {
		query.Limit(body.Length).Offset((body.BeginTableId - 1) * body.Length).Find(&keys)
	}
	//Mask the middle of each key before sending it to the frontend
	for i := range keys {
		keys[i].Key = maskKey(keys[i].Key)
	}
	//Return Result
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success", "data": keys})
	return
}

// maskKey keeps the first and last 5 characters, and replaces the middle with 8 asterisks
func maskKey(key string) string {
	if len(key) <= 10 {
		return "********"
	}
	return key[:5] + "********" + key[len(key)-5:]
}
