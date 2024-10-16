package crud_test

import (
	"context"
	"gotempl/controller"
	"gotempl/model"
	"gotempl/repository"
	"gotempl/service"
	"net/http"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Set up a test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	// Create test users
	testUsers := []model.User{
		{Uid: "1", Username: "user1", Role: "admin"},
		{Uid: "2", Username: "user2", Role: "user"},
	}
	for _, user := range testUsers {
		db.Create(&user)
	}

	// Set up repository and services
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := controller.NewUserHandler(userService)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Define routes
	r.GET("/admin/user", userHandler.UserCRUDHandler)
	r.POST("/api/user", userHandler.CreateUser)
	r.PUT("/api/user/:id", userHandler.UpdateUser)
	r.DELETE("/api/user/:id", userHandler.DeleteUser)

	return r, db
}

func TestUserCRUDHandler(t *testing.T) {
	r, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go server.ListenAndServe()
	defer server.Close()

	// Set up ChromeDP
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // Make the browser visible
		chromedp.Flag("start-maximized", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Navigate to the user CRUD page
	var content string
	err := chromedp.Run(ctx,
		chromedp.Navigate("http://localhost:8080/admin/user"),
		chromedp.OuterHTML("body", &content),
	)
	assert.NoError(t, err)
	assert.Contains(t, content, "User Management")

	// Test creating a new user
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("#uid", chromedp.ByID),
		chromedp.SetValue("#uid", "3", chromedp.ByID),
		chromedp.SetValue("#username", "newuser", chromedp.ByID),
		chromedp.SetValue("#role", "user", chromedp.ByID),
		chromedp.Click("#submitBtn"),
		chromedp.WaitVisible("#result .alert-success", chromedp.ByQuery),
	)
	assert.NoError(t, err)

	// Test editing a user
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("tr", chromedp.ByQuery),
		chromedp.SetValue("tr:last-child td:nth-child(2)", "updateduser", chromedp.ByQuery),
		chromedp.Click("tr:last-child button.btn-warning", chromedp.ByQuery),
		chromedp.WaitVisible("#result .alert-success", chromedp.ByQuery),
	)
	assert.NoError(t, err)

	// Test deleting a user
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("tr", chromedp.ByQuery),
		chromedp.Click("tr:last-child button.btn-danger", chromedp.ByQuery),
		//chromedp.WaitVisible("text/Are you sure you want to delete this user?", chromedp.BySearch),
		//chromedp.Click("text/OK", chromedp.BySearch),
		chromedp.WaitVisible("tr:last-child", chromedp.ByQuery),
	)
	assert.NoError(t, err)

	// Verify the number of users
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	assert.Equal(t, int64(2), userCount) // We added one and deleted one, so it should still be 2
}
