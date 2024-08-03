package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whoisaditya/golang-task-management-system/api/initializers"
	"github.com/whoisaditya/golang-task-management-system/api/models"
)

// CreateTask godoc
//
//	@Summary		Create an task
//	@Description	create task by json and bind by id
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	body		int	body	"Task params"
//	@Success		200	{string}	string		"ok"
//	@Failure		400	{object}	map[string]any
//	@Failure		404	{object}	map[string]any
//	@Failure		500	{object}	map[string]any
//	@Router			/task/create/	[post]
func CreateTask(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Find the user by ID
	var user models.User
	err := initializers.DB.First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	newTask := models.Task{
		Title:           body.Title,
		Description:     body.Description,
		CreatedBy:       userID,
		ActualStartTime: time.Now(),
		ActualEndTime:   time.Now(),
	}

	result := initializers.DB.Create(&newTask)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating task",
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully",
	})
}

func CreateTaskBulk(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Using dummy data to test
	// file, openErr := os.Open("data.csv")
	file_ptr, getErr := c.FormFile("taskBulkUpload")

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file",
		})
		return
	}

	file, openErr := file_ptr.Open()
	if openErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to open file",
		})
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)

	var tasks []models.Task

	// Skip the header row
	reader.Read()

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading CSV",
			})
			return
		}

		task := models.Task{
			Title:       row[0],
			Description: row[1],
			CreatedBy:   userID,
		}

		tasks = append(tasks, task)
	}

	// Bulk insert the tasks into the database
	insertErr := initializers.DB.Create(&tasks).Error
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error inserting tasks into the database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully uploaded %d tasks", len(tasks)),
	})
}

// Function to get All task
func GetTasks(c *gin.Context) {
	task_id := c.Query("task_id")

	var tasks []models.Task
	err := initializers.DB.Find(&tasks, task_id).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error getting tasks",
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, tasks)
}

func UpdateTask(c *gin.Context) {
	task_id := c.Query("task_id")

	// Find the user by ID
	var task models.Task
	err := initializers.DB.First(&task, task_id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Get Title and Description
	var body struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		ActualStartTime int64  `json:"actual_start_Time"`
		ActualEndTime   int64  `json:"actual_end_Time"`
		Seconds         int64  `json:"seconds"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields are empty",
		})
		return
	}

	if body.Title != "" {
		task.Title = body.Title
	}

	if body.Description != "" {
		task.Description = body.Description
	}

	if body.ActualStartTime != 0 {
		task.ActualStartTime = time.Unix(body.ActualStartTime, 0)
		task.ActualStartTime.Format(time.RFC3339)
	}

	if body.ActualEndTime != 0 {
		task.ActualEndTime = time.Unix(body.ActualEndTime, 0)
		task.ActualEndTime.Format(time.RFC3339)
	}

	if !task.ActualStartTime.Before(task.ActualEndTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Start Time cannot be before End Time",
		})
		return
	}

	err = initializers.DB.Save(&task).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error updating task",
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Details added successfully",
		"task":    task,
	})
}

func DeleteTask(c *gin.Context) {
	task_id := c.Query("task_id")

	// Find the user by ID
	var task models.Task
	err := initializers.DB.First(&task, task_id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	err = initializers.DB.Delete(&models.Task{}, task_id).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting task",
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
