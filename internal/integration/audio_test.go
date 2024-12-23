package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wahyuoi/sbc/internal/config"
	"github.com/wahyuoi/sbc/internal/event_handler"
	handler "github.com/wahyuoi/sbc/internal/http_handler"
	"github.com/wahyuoi/sbc/internal/middleware"
	"github.com/wahyuoi/sbc/internal/repository"
	"github.com/wahyuoi/sbc/internal/service"
)

func setupTestServer(t *testing.T) *gin.Engine {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file. ", err)
	}

	// Initialize repositories and services
	fileRepo := repository.NewFileRepositoryLocal()

	db, err := sql.Open("mysql", config.InitDB())
	require.NoError(t, err)
	t.Cleanup(func() {
		// todo: clean up table entries after each test.
		db.Close()
	})

	uow := repository.NewUnitOfWork(db, fileRepo)
	userService := service.NewUserService(uow)
	exerciseService := service.NewExerciseService(uow)
	audioService := service.NewAudioService()
	audioConverter := event_handler.NewAudioConverter(audioService, exerciseService)
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userHandler := handler.NewUserHandler(userService)
	exerciseHandler := handler.NewExerciseHandler(exerciseService, audioService, audioConverter)

	// Duplicate routes from main.go, it can be improved by using a shared router registration for both test and production..
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/audio/user/:user_id/phrase/:phrase_id", middleware.AuthMiddleware(), exerciseHandler.SubmitAudio)
	r.GET("/audio/user/:user_id/phrase/:phrase_id/:format", middleware.AuthMiddleware(), exerciseHandler.GetAudio)

	return r
}

func TestAudioSubmissionAndRetrieval(t *testing.T) {
	router := setupTestServer(t)

	// Step 1: Register a new user
	registerReq := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	registerBody, _ := json.Marshal(registerReq)
	registerResp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(registerResp, req)
	// Conflict if the test is re-executed using same database without cleaning up the table entries.
	assert.True(t, http.StatusCreated == registerResp.Code || http.StatusConflict == registerResp.Code)

	// Step 2: Login to get JWT token
	loginResp := httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(loginResp, req)
	assert.Equal(t, http.StatusOK, loginResp.Code)

	var loginResult map[string]string
	err := json.Unmarshal(loginResp.Body.Bytes(), &loginResult)
	require.NoError(t, err)
	token := loginResult["token"]
	require.NotEmpty(t, token)

	// Step 3: Submit audio file
	audioFile, err := os.Open("../../sample/three.m4a")
	require.NoError(t, err)
	defer audioFile.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("audio_file", filepath.Base(audioFile.Name()))
	require.NoError(t, err)
	_, err = io.Copy(part, audioFile)
	require.NoError(t, err)
	err = writer.Close()
	require.NoError(t, err)

	submitResp := httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/audio/user/2/phrase/1", &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(submitResp, req)
	assert.Equal(t, http.StatusCreated, submitResp.Code)

	// Step 4: Retrieve M4A format
	getM4aResp := httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/audio/user/2/phrase/1/m4a", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(getM4aResp, req)
	assert.Equal(t, http.StatusOK, getM4aResp.Code)
	assert.Equal(t, "audio/mp4", getM4aResp.Header().Get("Content-Type"))
	assert.NotEmpty(t, getM4aResp.Body.Bytes())
}

func TestAudioSubmissionErrors(t *testing.T) {
	router := setupTestServer(t)

	// Get JWT token first
	loginReq := map[string]string{
		"email":    "admin@example.com",
		"password": "admin123",
	}
	loginBody, _ := json.Marshal(loginReq)
	loginResp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(loginResp, req)
	require.Equal(t, http.StatusOK, loginResp.Code)

	var loginResult map[string]string
	err := json.Unmarshal(loginResp.Body.Bytes(), &loginResult)
	require.NoError(t, err)
	token := loginResult["token"]
	require.NotEmpty(t, token)

	tests := []struct {
		name           string
		userID         string
		phraseID       string
		setupAuth      bool
		expectedStatus int
	}{
		{
			name:           "missing auth token",
			userID:         "1",
			phraseID:       "1",
			setupAuth:      false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "accessing other user's resources",
			userID:         "999",
			phraseID:       "1",
			setupAuth:      true,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid phrase ID",
			userID:         "1",
			phraseID:       "999",
			setupAuth:      true,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			audioFile, err := os.Open("../../sample/three.m4a")
			require.NoError(t, err)
			defer audioFile.Close()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			part, err := writer.CreateFormFile("audio_file", filepath.Base(audioFile.Name()))
			require.NoError(t, err)
			_, err = io.Copy(part, audioFile)
			require.NoError(t, err)
			err = writer.Close()
			require.NoError(t, err)

			submitResp := httptest.NewRecorder()
			req = httptest.NewRequest("POST", fmt.Sprintf("/audio/user/%s/phrase/%s", tt.userID, tt.phraseID), &buf)
			if tt.setupAuth {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())
			router.ServeHTTP(submitResp, req)
			assert.Equal(t, tt.expectedStatus, submitResp.Code)
		})
	}
}
