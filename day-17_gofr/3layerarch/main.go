// @title           Task Management API
// @version         1.0
// @description     This is a simple API server for managing tasks and users.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Puneeth Gowda
// @contact.email  puneeth@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

package main

import (
	"3layerarch/migrations/migrations"
	"database/sql"
	"log"

	taskhandler "3layerarch/handler/task"
	userhandler "3layerarch/handler/user"

	taskservice "3layerarch/service/task"
	userservice "3layerarch/service/user"

	taskstore "3layerarch/store/task"
	userstore "3layerarch/store/user"

	_ "3layerarch/docs" // Swagger generated files

	_ "github.com/go-sql-driver/mysql"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()
	app.Migrate(migrations.All())

	// Initialize DB
	db, err := sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Println("Failed to connect to DB:", err)
		return
	}

	// Setup dependencies
	userStore := userstore.New(db)
	userService := userservice.New(userStore)
	userHandler := userhandler.New(userService)

	taskStore := taskstore.New(db)
	taskService := taskservice.New(taskStore, userService)
	taskHandler := taskhandler.New(taskService)

	// Register task routes
	app.POST("/task", taskHandler.CreateTask)
	app.GET("/task", taskHandler.ViewTasks)
	app.GET("/task/{id}", taskHandler.GetTask)
	app.PUT("/task/{id}", taskHandler.UpdateTask)
	app.DELETE("/task/{id}", taskHandler.DeleteTask)

	// Register user routes
	app.POST("/user", userHandler.CreateUser)
	app.GET("/user/{id}", userHandler.GetUser)

	app.Run()
}
