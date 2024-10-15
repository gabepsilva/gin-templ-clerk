package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestGetAllUsers(t *testing.T) {
	db, handler, router := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	router.GET("/user", handler.GetAllUsers)

	// Create some test users
	testUsers := []model.User{
		{Uid: "user1", Username: "testuser1", Role: "user"},
		{Uid: "user2", Username: "testuser2", Role: "admin"},
	}
	for _, user := range testUsers {
		db.Create(&user)
	}

	req, _ := http.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseUsers []model.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	assert.NoError(t, err)
	assert.Len(t, responseUsers, len(testUsers))

	// Check if the response contains all test users
	for _, testUser := range testUsers {
		found := false
		for _, responseUser := range responseUsers {
			if testUser.Uid == responseUser.Uid {
				assert.Equal(t, testUser.Username, responseUser.Username)
				assert.Equal(t, testUser.Role, responseUser.Role)
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("User with Uid %s not found in response", testUser.Uid))
	}
}

func TestGetUser(t *testing.T) {
	db, handler, router := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	router.GET("/user/:id", handler.GetUser)

	testUser := model.User{Uid: "testuser", Username: "testusername", Role: "user"}
	db.Create(&testUser)

	req, _ := http.NewRequest("GET", "/user/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser model.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.Uid, responseUser.Uid)
	assert.Equal(t, testUser.Username, responseUser.Username)
	assert.Equal(t, testUser.Role, responseUser.Role)

	// Test non-existent user
	req, _ = http.NewRequest("GET", "/user/nonexistent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateUser(t *testing.T) {
	db, handler, router := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	router.PUT("/user/:id", handler.UpdateUser)

	testUser := model.User{Uid: "testuser", Username: "testusername", Role: "user"}
	db.Create(&testUser)

	updatedUser := model.User{
		Uid:      testUser.Uid,
		Username: "updatedusername",
		Role:     "admin",
	}

	jsonUser, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/user/testuser", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser model.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Uid, responseUser.Uid)
	assert.Equal(t, updatedUser.Username, responseUser.Username)
	assert.Equal(t, updatedUser.Role, responseUser.Role)

	// Verify the user was actually updated in the database
	var dbUser model.User
	err = db.First(&dbUser, "uid = ?", testUser.Uid).Error
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Username, dbUser.Username)
	assert.Equal(t, updatedUser.Role, dbUser.Role)
}

func TestDeleteUser(t *testing.T) {
	db, handler, router := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	router.DELETE("/user/:id", handler.DeleteUser)

	testUser := model.User{Uid: "testuser", Username: "testusername", Role: "user"}
	db.Create(&testUser)

	req, _ := http.NewRequest("DELETE", "/user/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify the user was actually deleted from the database
	var dbUser model.User
	err := db.First(&dbUser, "uid = ?", testUser.Uid).Error
	assert.Error(t, err) // Should return an error as the user should not be found
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Test deleting a non-existent user
	req, _ = http.NewRequest("DELETE", "/user/nonexistent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code) // The API returns 204 even if the user doesn't exist
}
