package model

type User struct {
	Uid      string `json:"uid" form:"uid" gorm:"primaryKey;type:varchar(255)"`
	Username string `json:"username" form:"username"  gorm:"type:varchar(255);not null;unique"`
	Role     string `json:"role" form:"role" gorm:"type:varchar(50);not null;default:'user'"`
}
