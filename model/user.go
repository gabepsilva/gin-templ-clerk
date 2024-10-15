package model

type User struct {
	Uid      string `json:"uid" form:"uid" gorm:"primaryKey;type:varchar(255)" validate:"required,min=1"`
	Username string `json:"username" form:"username" gorm:"type:varchar(255);not null;unique" validate:"required,min=1"`
	Role     string `json:"role" form:"role" gorm:"type:varchar(50);not null;default:'user'" validate:"oneof=user admin"`
}
