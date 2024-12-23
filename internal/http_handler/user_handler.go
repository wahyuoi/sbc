package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/wahyuoi/sbc/internal/common"
	"github.com/wahyuoi/sbc/internal/config"
	"github.com/wahyuoi/sbc/internal/model"
	"github.com/wahyuoi/sbc/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// todo: do email and password validation
	// currently, we skip as it is not the focus of this task.

	if err := h.userService.Register(c.Request.Context(), req.Email, req.Password); err != nil {
		log.Println(err)
		status, message := common.ConvertErrorToHTTPStatus(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// todo: do email validation

	user, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{Token: tokenString})
}
