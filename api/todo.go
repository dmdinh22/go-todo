package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type (
	TodoModel struct {
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

func CreateTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))

	if c.PostForm("completed") == "true" {
		completed = 1
	} else if c.PostForm("completed") == "false" {
		completed = 0
	}

	todo := TodoModel{Title: c.PostForm("title"), Message: c.PostForm("message"), Completed: completed}

	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item successfully created.",
		"resourceId": todo.ID,
	})
}

func GetAllTodo(c *gin.Context) {
	var todos []TodoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, noTodoFoundError)
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
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   _todos,
	})
}

func GetSingleTodo(c *gin.Context) {
	var todo TodoModel
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

func UpdateTodo(c *gin.Context) {
	var todo TodoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, noTodoFoundError)
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

func DeleteTodo(c *gin.Context) {
	var todo TodoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, noTodoFoundError)
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully.",
	})
}

var noTodoFoundError = gin.H{
	"status":  http.StatusNotFound,
	"message": "No todo found!",
}
