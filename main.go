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
		v1.GET("/", getAllTodo)
		v1.GET("/:id", getSingleTodo)
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
		Message    string `json:"message"`
		Completed  int    `json:"completed"`
	}

	// formatted todo
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Message   string `json:"message"`
		Completed bool   `json:"completed`
	}
)

func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Message: c.PostForm("message"), Completed: completed}

	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item successfully created.",
		"resourceId": todo.ID,
	})
}

func getAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No todo found!",
		})
		return
	}

	// transforms todos to build resp
	for _, item := range todos {
		completed := false

		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}

		_todos = append(_todos, transformedTodo{
			ID:        item.ID,
			Title:     item.Title,
			Message:   item.Message,
			Completed: completed,
		})

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   _todos,
		})
	}
}

func getSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No todo found.",
		})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{
		ID:        todo.ID,
		Title:     todo.Title,
		Message:   todo.Message,
		Completed: completed,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   _todo,
	})
}

func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No todo found!",
		})
		return
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	db.Model(&todo).Update("message", c.PostForm("message"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo updated successfully.",
	})
}
