package main

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)


type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// request structure
type CreateTodoRequest struct {
	Title string `json:"title"`
}

var todos = []Todo{}
var mu sync.Mutex

func main() {
	r := gin.Default()

	// get to get all todos
	r.GET("/todos", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		c.JSON(200, todos)
	})

	//post to add todo item
	r.POST("/todos", func(c *gin.Context) {
		var req CreateTodoRequest

		if err := c.BindJSON(&req); err != nil {
			c.Status(400)
			return
		}

		//locking the todos slice for concurrent access
		mu.Lock()
		defer mu.Unlock()

		todo := Todo{
			ID:    len(todos) + 1,
			Title: req.Title,
			Done:  false,
		}

		todos = append(todos, todo)

		c.JSON(201, todo)
	})

	// put to mark todo as done
	r.PUT("/todos", func(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.Status(400)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Status(400)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Done = true
			c.JSON(200, todos[i])
			return
		}
	}

	c.Status(404)
   })

   // delete to remove todo item
   r.DELETE("/todos", func(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.Status(400)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Status(400)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.Status(204)
			return
		}
	}

	c.Status(404)
	})

	//start the server
	r.Run(":8080")
}
