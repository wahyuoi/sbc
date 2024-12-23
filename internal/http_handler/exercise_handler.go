package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wahyuoi/sbc/internal/common"
	"github.com/wahyuoi/sbc/internal/event_handler"
	"github.com/wahyuoi/sbc/internal/model"
	"github.com/wahyuoi/sbc/internal/service"
)

type ExerciseHandler struct {
	exerciseService *service.ExerciseService
	audioService    *service.AudioService
	audioConverter  *event_handler.AudioConverter
}

func NewExerciseHandler(
	exerciseService *service.ExerciseService,
	audioService *service.AudioService,
	audioConverter *event_handler.AudioConverter,
) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
		audioService:    audioService,
		audioConverter:  audioConverter,
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
	// But for now, we just do it this way to make it simple.
	fileExt := filepath.Ext(audioFileHeader.Filename)
	fileFormat, err := model.GetAudioFormatType(fileExt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio file format"})
		return
	}
	// This is to limit the file format for upload to only M4A.
	// If we want to support more file format, we can add it to audioProps in model/audio.go.
	if !fileFormat.IsForUpload() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio file format"})
		return
	}

	// todo: we also need to check the file size, and audio duration.
	// options: using ffprobe (https://trac.ffmpeg.org/wiki/FFprobeTips) which comes with ffmpeg.
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
		status, message := common.ConvertErrorToHTTPStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	// Converting audio file to other format in background, so it won't block the response.
	// We could do this because the requirement says that the converted audio file is only stored in the database.
	// We could use a queue to do this, but for now we just do it in a goroutine.
	go func() {
		// using new context because the main context will be closed when the response is sent, and could cause the audio converter cancelled.
		ctx := context.Background()
		err := h.audioConverter.ConvertAudio(ctx, userID, phraseID, model.AudioFormatTypeM4a, model.AudioFormatTypeWav)
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
	audioFormat, err := model.GetAudioFormatType(c.Param("format"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio file format"})
		return
	}
	if !audioFormat.IsForDownload() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio file format"})
		return
	}

	audio, err := h.exerciseService.GetAudio(c.Request.Context(), userID, phraseID, audioFormat)
	if err != nil {
		log.Println(err)
		status, message := common.ConvertErrorToHTTPStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.Data(http.StatusOK, audioFormat.GetMimeType(), audio)
}
