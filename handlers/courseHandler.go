package handlers

import (
	"student-enrollment-system/models"
	"net/http"
	"strconv"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
)
var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
}

// Handlers
func ListCourses(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description, credits FROM courses")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve courses: %v", err)
		return
	}
	defer rows.Close()

	var courses []models.Course

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Credits); err != nil {
			c.String(http.StatusInternalServerError, "Failed to scan course: %v", err)
			return
		}
		courses = append(courses, course)
	}

	c.HTML(http.StatusOK, "courses.html", gin.H{"courses": courses})
}

func AddCourse(c *gin.Context) {
    courseName := c.PostForm("name") 
    courseDescription := c.PostForm("description") 
    courseCreditsStr := c.PostForm("credit") 

    if courseCreditsStr == "" {
        c.String(http.StatusBadRequest, "Course credit is required and cannot be empty")
        return
    }

    courseCredits, err := strconv.Atoi(courseCreditsStr)
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid course credit value, must be a number: %v", err)
        return
    }

    _, err = db.Exec("INSERT INTO courses (name, description, credit) VALUES (?, ?, ?)",
        courseName, courseDescription, courseCredits)

    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to add course: %v", err)
        return
    }

    c.Redirect(http.StatusSeeOther, "/courses")
}

func getCourseForEdit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	var course models.Course
	err := db.QueryRow("SELECT id, name, description, credits FROM courses WHERE id = ?", id).
		Scan(&course.ID, &course.Name, &course.Description, &course.Credits)

	if err != nil {
		// Log the error for debugging
		fmt.Println("Error fetching course: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	fmt.Println(course)

	c.HTML(http.StatusOK, "editCourse.html", gin.H{
		"course": course,
	})
}

func updateCourse(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	description := c.PostForm("description")
	creditsStr := c.PostForm("credits")

	if creditsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credits cannot be empty"})
		return
	}

	credits, err := strconv.Atoi(creditsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credits format"})
		return
	}

	result, err := db.Exec("UPDATE courses SET name = ?, description = ?, credits = ? WHERE id = ?", name, description, credits, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found or not updated"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/courses")
}


func deleteCourse(c *gin.Context) {
	id := c.PostForm("id")
	_, _ = db.Exec("DELETE FROM courses WHERE id = ?", id)
	c.Redirect(http.StatusSeeOther, "/courses")
}