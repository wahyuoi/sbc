package handler

import (
	"context"
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
	audioService    service.AudioService
}

func NewExerciseHandler(
	exerciseService service.ExerciseService,
	audioService service.AudioService,
) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
		audioService:    audioService,
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
		return
	}

	// Converting audio file to other format in background, so it won't block the response.
	// We could do this because the requirement says that the converted audio file is only stored in the database.
	// We could use a queue to do this, but for now we just do it in a goroutine.
	go func() {
		ctx := context.Background()
		originalAudioBytes, err := h.exerciseService.GetAudio(ctx, userID, phraseID, model.AudioFormatTypeM4a)
		if err != nil {
			log.Println(err)
			return
		}
		newAudioBytes, err := h.audioService.ConvertAudio(ctx, originalAudioBytes, model.AudioFormatTypeWav)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.exerciseService.SubmitAudio(ctx, userID, phraseID, newAudioBytes, model.AudioFormatTypeWav)
		if err != nil {
			log.Println(err)
		}
	}()

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
