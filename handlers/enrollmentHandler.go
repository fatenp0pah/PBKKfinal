package handlers

import (
	"net/http"
	"student-enrollment-system/models"
	"github.com/gin-gonic/gin"
	"database/sql"
    "strconv"

)

func ListEnrollments(c *gin.Context, db *sql.DB) {
    enrollments, err := models.GetEnrollments(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query enrollments"})
        return
    }

    // Get the list of students and courses
    students, err := models.GetAllStudents(db) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
        return
    }

    courses, err := models.FetchAllCourses(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
        return
    }

    // Pass the data to the template
    c.HTML(http.StatusOK, "enrollments.html", gin.H{
        "enrollments": enrollments,
        "students":    students,
        "courses":     courses,
    })
}
func ShowEnrollmentsPage(c *gin.Context, db *sql.DB) {
    // Fetch students from the database
    students, err := models.GetAllStudents(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load students"})
        return
    }

    // Fetch courses from the database
    courses, err := models.FetchAllCourses(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load courses"})
        return
    }

    // Fetch enrollment records
    enrollments, err := models.GetEnrollments(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load enrollments"})
        return
    }

    // Render the template with the data
    c.HTML(http.StatusOK, "enrollments.html", gin.H{
        "students": students,
        "courses": courses,
        "enrollments": enrollments,
    })
}



func AssignEnrollment(c *gin.Context, db *sql.DB) {
    // Parse form values
    studentIDStr := c.PostForm("student_id")
    courseIDStr := c.PostForm("course_id")

    // Convert studentID to an integer
    studentID, err := strconv.Atoi(studentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
        return
    }

    // Convert courseID to an integer
    courseID, err := strconv.Atoi(courseIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
        return
    }

    // Insert enrollment into the database
    err = models.CreateEnrollment(db, studentID, courseID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign enrollment"})
        return
    }

    // Redirect back to the enrollment page
    c.Redirect(http.StatusFound, "/enrollments")
}
