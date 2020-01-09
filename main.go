package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}

	router.Run(":5000")
}

func init() {
	//open db connection
	var err error
	db, err = gorm.Open("mysql", "root:Password1@/demo?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	}

	// migrate schema
	db.AutoMigrate(&todoModel{})
}

type (
	todoModel struct {
		gorm.Model        //generates a model struct for ID, CreatedAt, UpdatedAt, DeletedAt
		Title      string `json:"title"`
		Completed  int    `json:"completed"`
	}

	// formatted todo
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed`
	}
)

func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}

	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item successfully created.", "resourceId": todo.ID})
}
