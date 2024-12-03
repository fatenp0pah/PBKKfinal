package handlers

import (
	"student-enrollment-system/controllers"
	"student-enrollment-system/models"
	"net/http"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
)

func ListStudentsHandler(c *gin.Context, db *sql.DB) {
	students, err := controllers.FetchAllStudents(db)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching students: %v", err)
		return
	}
	c.HTML(http.StatusOK, "students.html", gin.H{"students": students})
}

func AddStudentHandler(c *gin.Context, db *sql.DB) {
	// Get form values from POST request
	name := c.PostForm("name")
	age, _ := strconv.Atoi(c.PostForm("age"))
	phone := c.PostForm("phone")
	image := c.PostForm("image")

	// Create a new student model
	student := models.Student{Name: name, Age: age, Phone: phone, Image: image}

	// Call the controller to add the student
	err := controllers.AddStudent(db, student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add student"})
		return
	}

	// Redirect to the student list after successful addition
	c.Redirect(http.StatusSeeOther, "/students")
}

func GetStudentForEditHandler(c *gin.Context, db *sql.DB) {
	// Get the student ID from the URL parameter
	id := c.Param("id")
	student, err := controllers.FetchStudentByID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.HTML(http.StatusOK, "editStudent.html", gin.H{"student": student})
}

func UpdateStudentHandler(c *gin.Context, db *sql.DB) {
	// Get the form values from the POST request
	id, _ := strconv.Atoi(c.PostForm("id"))
	name := c.PostForm("name")
	age, _ := strconv.Atoi(c.PostForm("age"))
	phone := c.PostForm("phone")
	image := c.PostForm("image")

	// Create a student object for updating
	student := models.Student{ID: id, Name: name, Age: age, Phone: phone, Image: image}
	_, err := controllers.ModifyStudent(db, student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	// Redirect to the student list after successful update
	c.Redirect(http.StatusSeeOther, "/students")
}

func DeleteStudentHandler(c *gin.Context, db *sql.DB) {
	// Get the student ID from the URL parameter or form data
	id := c.DefaultPostForm("id", "") // DefaultPostForm gives an empty string if "id" is missing.

	if id == "" {
		// If the id is empty, return an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}

	// Call the controller to remove the student using the given ID
	err := controllers.DeleteStudent(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	// Redirect to the student list page after successful deletion
	c.Redirect(http.StatusSeeOther, "/students")
}
