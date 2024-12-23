package main

import (
	"database/sql"
	"log"

	"github.com/wahyuoi/sbc/internal/config"
	"github.com/wahyuoi/sbc/internal/event_handler"
	handler "github.com/wahyuoi/sbc/internal/http_handler"
	"github.com/wahyuoi/sbc/internal/middleware"
	"github.com/wahyuoi/sbc/internal/repository"
	"github.com/wahyuoi/sbc/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadEnv()

	db, err := sql.Open("mysql", config.InitDB())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	fileRepository := repository.NewFileRepositoryLocal()
	uow := repository.NewUnitOfWork(db, fileRepository)
	userService := service.NewUserService(uow)
	exerciseService := service.NewExerciseService(uow)
	audioService := service.NewAudioService()
	audioConverter := event_handler.NewAudioConverter(audioService, exerciseService)
	userHandler := handler.NewUserHandler(userService)
	exerciseHandler := handler.NewExerciseHandler(exerciseService, audioService, audioConverter)

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	r.POST("/audio/user/:user_id/phrase/:phrase_id", middleware.AuthMiddleware(), exerciseHandler.SubmitAudio)
	r.GET("/audio/user/:user_id/phrase/:phrase_id/:format", middleware.AuthMiddleware(), exerciseHandler.GetAudio)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
