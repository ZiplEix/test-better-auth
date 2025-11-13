package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ZiplEix/test-better-auth/backend/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Todo struct {
	ID        int64  `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// GET /todos
func GetTodos(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	rows, err := database.DB.Query(context.Background(),
		"SELECT id, user_id, title, completed FROM todos WHERE user_id=$1 ORDER BY id DESC", userID)
	if err != nil {
		return fmt.Errorf("")
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Completed)
		todos = append(todos, t)
	}

	return c.JSON(http.StatusOK, todos)
}

// POST /todos
func CreateTodo(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	var body struct {
		Title string `json:"title"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	_, err := database.DB.Exec(context.Background(),
		"INSERT INTO todos(user_id, title) VALUES($1,$2)", userID, body.Title)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}

// Delete todo (DELETE /todos/:id)
func DeleteTodo(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	cmd := `
        DELETE FROM todos
        WHERE id = $1 AND user_id = $2
    `
	res, err := database.DB.Exec(context.Background(), cmd, id, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "todo not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// Toggle todo completed (PATCH /todos/:id/toggle)
func ToggleTodo(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	// simple toggle: flip completed for this user's todo
	cmd := `
        UPDATE todos
        SET completed = NOT completed
        WHERE id = $1 AND user_id = $2
    `
	res, err := database.DB.Exec(context.Background(), cmd, id, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "todo not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "toggled"})
}
