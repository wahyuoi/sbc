package handler

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wahyuoi/sbc/internal/model"
	"github.com/wahyuoi/sbc/internal/service"
)

type ExerciseHandler struct {
	exerciseService service.ExerciseService
}

func NewExerciseHandler(
	exerciseService service.ExerciseService,
) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
	}
}

func (h *ExerciseHandler) SubmitAudio(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	loggedInUserID := c.GetInt64("user_id")
	if loggedInUserID != int64(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to submit audio for this user"})
		return
	}
	phraseID, err := strconv.Atoi(c.Param("phrase_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phrase ID"})
		return
	}
	audioFileHeader, err := c.FormFile("audio_file")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get audio file"})
		return
	}

	// I think this is not the correct way to get the file format, as we could edit file extension manually.
	// We should use the file header to get the file format, or use existing library like https://github.com/gabriel-vasile/mimetype
	fileExt := filepath.Ext(audioFileHeader.Filename)
	fileFormat := model.AudioFormatType(strings.TrimPrefix(fileExt, "."))
	if fileFormat != model.AudioFormatTypeWav && fileFormat != model.AudioFormatTypeM4a {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio file format"})
		return
	}

	audioFile, err := audioFileHeader.Open()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open audio file"})
		return
	}
	defer audioFile.Close()

	audioBytes, err := io.ReadAll(audioFile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read audio file"})
		return
	}

	err = h.exerciseService.SubmitAudio(c.Request.Context(), userID, phraseID, audioBytes, fileFormat)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit audio"})
		// TODO: do we need to delete the audio file if we failed to record it in the database?
		return
	}

	// TODO: convert audio file to wav

	c.JSON(http.StatusCreated, gin.H{"message": "Audio submitted successfully"})
}

func (h *ExerciseHandler) GetAudio(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	loggedInUserID := c.GetInt64("user_id")
	if loggedInUserID != int64(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to submit audio for this user"})
		return
	}
	phraseID, err := strconv.Atoi(c.Param("phrase_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phrase ID"})
		return
	}
	audioFormat := model.AudioFormatType(c.Param("format"))
	if audioFormat != model.AudioFormatTypeWav && audioFormat != model.AudioFormatTypeM4a {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio file format"})
		return
	}

	audio, err := h.exerciseService.GetAudio(c.Request.Context(), userID, phraseID, audioFormat)
	if err != nil {
		log.Println(err)
		// TODO: handle record not found
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get audio"})
		return
	}

	c.Data(http.StatusOK, audioFormat.GetMimeType(), audio)
}
