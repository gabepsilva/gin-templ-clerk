package repository

import (
	"gotempl/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, "uid = ?", id).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.DB.Delete(&model.User{}, "uid = ?", id).Error
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
