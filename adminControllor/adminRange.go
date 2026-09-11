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
	BeginTableId    int                  `json:"begin_table_id"`
	Length          int                  `json:"length"`
	SearchCondition AdminSearchCondition `json:"search_condition"`
}

type AdminSearchCondition struct {
	//If The Field is Empty, It Means Ignore This Condition
	Username string `json:"username"`
	UUHash   string `json:"uu_hash"`

	//Permission Filter (Bool, false means ignore)
	CanAddAdmin         bool `json:"can_add_admin"`
	CanDeleteAdmin      bool `json:"can_delete_admin"`
	CanEditAdmin        bool `json:"can_edit_admin"`
	CanGetAdmin         bool `json:"can_get_admin"`
	CanOperateUser      bool `json:"can_operate_user"`
	CanOperateCharacter bool `json:"can_operate_character"`
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
	//Build Search Query (only non-empty fields, combined with AND)
	condition := body.SearchCondition
	var conditions []string
	var args []interface{}
	if condition.Username != "" {
		conditions = append(conditions, "username LIKE ?")
		args = append(args, "%"+condition.Username+"%")
	}
	if condition.UUHash != "" {
		conditions = append(conditions, "uu_hash LIKE ?")
		args = append(args, "%"+condition.UUHash+"%")
	}
	//Permission Filter: only checked (true) fields are applied, all combined with AND
	if condition.CanAddAdmin {
		conditions = append(conditions, "permission ->> 'can_add_admin' = 'true'")
	}
	if condition.CanDeleteAdmin {
		conditions = append(conditions, "permission ->> 'can_delete_admin' = 'true'")
	}
	if condition.CanEditAdmin {
		conditions = append(conditions, "permission ->> 'can_edit_admin' = 'true'")
	}
	if condition.CanGetAdmin {
		conditions = append(conditions, "permission ->> 'can_get_admin' = 'true'")
	}
	if condition.CanOperateUser {
		conditions = append(conditions, "permission ->> 'can_operate_user' = 'true'")
	}
	if condition.CanOperateCharacter {
		conditions = append(conditions, "permission ->> 'can_operate_character' = 'true'")
	}
	query := db.Model(&sqlTable.Admin{})
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	//Must Be Ordered, Otherwise Limit/Offset Follows Physical Row Order And Pages Are Not Stable
	query = query.Order("id ASC")
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
