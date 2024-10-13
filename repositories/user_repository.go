package repositories

import (
	"gotempl/database"
	"gotempl/models"
)

type UserRepository struct{}

func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return database.DB.Delete(&models.User{}, "uid = ?", id).Error
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
