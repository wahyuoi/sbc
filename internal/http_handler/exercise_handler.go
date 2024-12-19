package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wahyuoi/sbc/internal/service"
)

type ExerciseHandler struct {
	exerciseService service.ExerciseService
}

func NewExerciseHandler(exerciseService service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{exerciseService: exerciseService}
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

	audio, err := c.FormFile("audio_file")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get audio file"})
		return
	}

	audioPath := "dummyPath"

	err = h.exerciseService.SubmitAudio(c.Request.Context(), userID, phraseID, audioPath, audio.Filename)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit audio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audio submitted successfully"})
}
