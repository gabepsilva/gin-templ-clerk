package controller

import (
	"bytes"
	"encoding/json"
	"gotempl/controller/service"
	"gotempl/model"
	"gotempl/repositories"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestEnvironment(t *testing.T) (*gorm.DB, *UserHandler, *gin.Engine) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	repo := repositories.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := NewUserHandler(service)

	router := gin.Default()
	return db, handler, router
}

func TestCreateUser(t *testing.T) {
	db, handler, router := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	router.POST("/user", handler.CreateUser)

	t.Run("Valid user creation", func(t *testing.T) {
		user := model.User{
			Uid:      "test-uid",
			Username: "testuser",
			Role:     "user",
		}

		jsonUser, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var responseUser model.User
		err := json.Unmarshal(w.Body.Bytes(), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Uid, responseUser.Uid)
		assert.Equal(t, user.Username, responseUser.Username)
		assert.Equal(t, user.Role, responseUser.Role)

		// Verify the user was actually created in the database
		var dbUser model.User
		err = db.First(&dbUser, "uid = ?", user.Uid).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Username, dbUser.Username)
		assert.Equal(t, user.Role, dbUser.Role)
	})

	t.Run("Invalid user creation", func(t *testing.T) {
		invalidUser := struct {
			InvalidField string `json:"invalid_field"`
		}{
			InvalidField: "test",
		}

		jsonUser, _ := json.Marshal(invalidUser)
		req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
