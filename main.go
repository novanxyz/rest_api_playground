package main

import (
	"fmt"
	"net/http"
	"novanxyz/api"
	"novanxyz/api/controller"
	"novanxyz/config"
	_ "novanxyz/docs"
	"novanxyz/models"
	"novanxyz/repository"
	"novanxyz/service"
	"novanxyz/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Application struct {
	DB     *gorm.DB
	Engine *gin.Engine
}

func newApplication(db *gorm.DB) *Application {

	validate := utils.CreateValidator()

	db.Table("tasks").AutoMigrate(&models.Task{})
	db.Table("task_files").AutoMigrate(&models.TaskFile{})

	taskRepository := repository.NewTaskRepository(db)
	taskFileRepository := repository.NewTaskFileRepository(db)

	taskService := service.NewTaskService(taskRepository, taskFileRepository, validate)
	taskController := controller.NewTaskController(taskService)
	engine := api.NewRouter(taskController)

	return &Application{
		DB:     db,
		Engine: engine,
	}
}

// @title 	Task Service API
// @version	1.0
// @description A Task service API in Go using Gin framework

// @host 	localhost:9001
// @BasePath /api
func main() {
	port := utils.Getenv("PORT", "9001")
	log.Info().Msg(fmt.Sprintf("Started Server : %s:%s !", utils.Getenv("DB_NAME", ""), port))
	db := config.DatabaseConnection()

	app := newApplication(db)
	server := &http.Server{
<<<<<<< HEAD
		Addr:    fmt.Sprintf("%s:%s", "127.0.0.1", port),
=======
		Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", port),
>>>>>>> 3252580ef0a5dfa3eb455fc35a757b43233beb6d
		Handler: app.Engine,
	}

	err := server.ListenAndServe()
	utils.ErrorPanic(err)
}
