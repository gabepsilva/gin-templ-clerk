package service

import (
	"errors"
	"gotempl/model"
	"gotempl/repositories"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repo     *repositories.UserRepository
	validate *validator.Validate
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *UserService) CreateUser(user *model.User) error {
	if err := s.validate.Struct(user); err != nil {
		return err
	}

	if user.Uid == "" || user.Username == "" {
		return errors.New("uid and username are required")
	}

	return s.repo.Create(user)
}

func (s *UserService) GetUser(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) UpdateUser(user *model.User) error {
	if err := s.validate.Struct(user); err != nil {
		return err
	}
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

// Additional method to match the handler
func (s *UserService) GetAllUser() ([]model.User, error) {
	return s.GetAllUsers()
}
