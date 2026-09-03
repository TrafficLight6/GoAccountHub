package sqlTable

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"column:username"` //CONSTANT!!
	PasswordHash string `gorm:"column:password_hash"`
	UUHash       string `gorm:"column:uu_hash"` //CONSTANT!!

	Permission Permission `gorm:"column:permission;type:jsonb;serializer:json"`
}

type Permission struct {
	//All Fields must be Boolean Type
	//Admin Permission
	CanAddAdmin    bool `gorm:"column:can_add_admin" json:"can_add_admin"`
	CanDeleteAdmin bool `gorm:"column:can_delete_admin" json:"can_delete_admin"`
	CanEditAdmin   bool `gorm:"column:can_edit_admin" json:"can_edit_admin"`
	CanGetAdmin    bool `gorm:"column:can_get_admin" json:"can_get_admin"`
	//User Permission
	CanOperateUser bool `gorm:"column:can_operate_user" json:"can_operate_user"`
	//Character Permission
	CanOperateCharacter bool `gorm:"column:can_operate_character" json:"can_operate_character"`
}

type AdminForEdit struct {
	Password   string     `gorm:"column:password"`
	Permission Permission `gorm:"column:permission;type:jsonb;serializer:json"`
}
