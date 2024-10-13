package models

type User struct {
	Uid      string `json:"uid" gorm:"primaryKey;type:varchar(255)"`
	Username string `json:"username" gorm:"type:varchar(255);not null;unique"`
	Role     string `json:"role" gorm:"type:varchar(50);not null;default:'user'"`
}
