package todo

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	store := NewTodoStore(db)
	svc := NewTodoService(store)
	handler := NewTodoHandler(svc)
	r.GET("/todos/", handler.GetTodos)
	r.GET("/todos/:id", handler.GetTodo)
	r.POST("/todos/", handler.CreateTodo)
	r.PUT("/todos/:id", handler.ToggleTodoState) // Toggle done state
	r.PATCH("/todos/:id", handler.UpdateTodo)    // Update all status
	r.DELETE("/todos/:id", handler.DestroyTodo)
}
