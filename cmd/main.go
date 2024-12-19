package main

import (
	"database/sql"
	"log"

	"github.com/wahyuoi/sbc/internal/config"
	handler "github.com/wahyuoi/sbc/internal/http_handler"
	"github.com/wahyuoi/sbc/internal/middleware"
	"github.com/wahyuoi/sbc/internal/repository"
	"github.com/wahyuoi/sbc/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", config.InitDB())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	uow := repository.NewSqlUnitOfWork(db)
	userService := service.NewUserService(uow)
	userHandler := handler.NewUserHandler(userService)
	exerciseService := service.NewExerciseService(uow)
	exerciseHandler := handler.NewExerciseHandler(exerciseService)

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	r.POST("/audio/user/:user_id/phrase/:phrase_id", middleware.AuthMiddleware(), exerciseHandler.SubmitAudio)

	// temporary, just for testing auth
	r.GET("/hello", middleware.AuthMiddleware(), userHandler.Hello)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
