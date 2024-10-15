package controller

import (
	"errors"
	"gotempl/model"
	"gotempl/repositories"
	"gotempl/views/crud"
	"gotempl/views/layout"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repo repositories.UserRepository
}

/*
func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}
*/

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with the provided information
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      model.User  true  "User information"
// @Success      201   {object}  model.User
// @Failure      400   {object}  object
// @Failure      500   {object}  object
// @Router       /user [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Retrieve a list of all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.User
// @Failure      500  {object}  map[string]string
// @Router       /user [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get a user by ID
// @Description  Retrieve a user's information using their ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  model.User
// @Failure      400  {object}  object
// @Failure      404  {object}  object
// @Router       /user/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update a user's information in the system
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id    path      string     true  "User ID"
// @Param        user  body      model.User true  "Updated user information"
// @Success      200   {object}  model.User
// @Failure      400   {object}  object
// @Failure      500   {object}  object
// @Router       /user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Uid = string(id)
	if err := h.Repo.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user from the system using their ID.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  object
// @Router       /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UserCRUDHandler godoc
// @Summary      This is a non-REST endpoint that returns an HTML page - not JSON data
// @Description  Fetches all users and renders an HTML page with a CRUD form for user management (non-REST endpoint)
// @Tags         User
// @Produce      html
// @Success      200  {string}  string  "HTML page content"
// @Router       /admin/user/crud [get]
// @Notes
func (h *UserHandler) UserCRUDHandler(c *gin.Context) {
	users, err := h.Repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	//print(users)
	layout.Render(c, 200, crud.UserForm(users, &model.User{}))
}
