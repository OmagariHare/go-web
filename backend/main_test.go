package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-web/config"
	"go-web/database"
	"go-web/dtos"
	"go-web/models"
	"go-web/routers"
	"go-web/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	testDB  *gorm.DB
	testCfg *config.Config
	router  *gin.Engine
)

func TestMain(m *testing.M) {
	// Setup
	gin.SetMode(gin.TestMode)
	testCfg = setupTestConfig()
	utils.InitLogger(testCfg.Log.Level, "", 0, 0, 0, false) // Log to console for tests
	testDB = setupTestDatabase(testCfg)
	router = routers.SetupRouter(testCfg)

	// Run tests
	exitCode := m.Run()

	// Teardown
	teardownTestDatabase()

	// Exit
	os.Exit(exitCode)
}

func setupTestConfig() *config.Config {
	// Use a dedicated test configuration
	return &config.Config{
		App: config.AppConfig{
			DefaultRole: "user",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "mysecretpassword",
			DBName:   "go_web_test",
			SSLMode:  "disable",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: 3600,
		},
		Server: config.ServerConfig{
			Port:           8081, // Use a different port for testing
			AllowedOrigins: []string{"*"},
		},
		Casbin: config.CasbinConfig{
			Model: `[request_definition]
    r = sub, obj, act
    
    [policy_definition]
    p = sub, obj, act
    
    [role_definition]
    g = _, _
    
    [policy_effect]
    e = some(where (p.eft == allow))
    
    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`,
		},
		Log: config.LogConfig{
			Level: "debug",
		},
		RateLimiter: config.RateLimiterConfig{
			Period: "1m",
			Limit:  10,
		},
	}
}

func setupTestDatabase(cfg *config.Config) *gorm.DB {
	// Connect to the PostgreSQL server
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Port, cfg.Database.SSLMode)
	db, err := gorm.Open(database.NewPostgresDialector(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to postgres server: %v", err))
	}

	// Create the test database
	db.Exec("DROP DATABASE IF EXISTS " + cfg.Database.DBName)
	db.Exec("CREATE DATABASE " + cfg.Database.DBName)

	// Connect to the newly created test database
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)
	testDB, err = gorm.Open(database.NewPostgresDialector(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}

	// Run migrations
	err = testDB.AutoMigrate(&models.User{}, &models.Role{}, &gormadapter.CasbinRule{})
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	// Seed roles
	seedRoles(testDB)

	database.DB = testDB // Override the global DB instance for testing
	return testDB
}

func seedRoles(db *gorm.DB) {
	roles := []models.Role{{Name: "admin"}, {Name: "user"}}
	for _, role := range roles {
		db.Create(&role)
	}
}

func teardownTestDatabase() {
	// Drop the test database
	sqlDB, _ := testDB.DB()
	err := sqlDB.Close()
	if err != nil {
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=%s",
		testCfg.Database.Host, testCfg.Database.User, testCfg.Database.Password, testCfg.Database.Port, testCfg.Database.SSLMode)
	dbConn, err := gorm.Open(database.NewPostgresDialector(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to postgres server for teardown: %v", err))
	}
	dbConn.Exec("DROP DATABASE " + testCfg.Database.DBName)
}

// clearTables clears data from tables that are modified during tests.
func clearTables(db *gorm.DB) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&gormadapter.CasbinRule{})
}

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestUserRegistration(t *testing.T) {
	defer clearTables(testDB)

	w := httptest.NewRecorder()
	userPayload := gin.H{
		"username": "testuser",
		"password": "password123",
		"email":    "test@example.com",
	}
	body, _ := json.Marshal(userPayload)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dtos.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Token, "Token should not be empty")
	assert.Equal(t, "testuser", response.User.Username, "Username should match")
}

func TestUserLogin(t *testing.T) {
	defer clearTables(testDB)

	// First, register a user
	regUser := gin.H{
		"username": "loginuser",
		"password": "password123",
		"email":    "login@example.com",
	}
	regBody, _ := json.Marshal(regUser)
	regReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(regBody))
	regReq.Header.Set("Content-Type", "application/json")

	regRecorder := httptest.NewRecorder()
	router.ServeHTTP(regRecorder, regReq)
	assert.Equal(t, http.StatusCreated, regRecorder.Code)

	// Then, attempt to log in
	w := httptest.NewRecorder()
	loginCreds := gin.H{
		"username": "loginuser",
		"password": "password123",
	}
	body, _ := json.Marshal(loginCreds)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.NotEmpty(t, response["token"])
}
