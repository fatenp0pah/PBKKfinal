package main

import (
	"database/sql"
	"log"
	"net/http"
	"student-enrollment-system/controllers"
	"student-enrollment-system/models"
	"student-enrollment-system/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main() {
	// Open database connection
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/student_enrollment_system")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Start Gin web server
	r := gin.Default()

	// Serve static files for images, CSS, JS, etc.
	r.Static("/assets", "./assets")
	r.Static("/static", "./static")

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve `index.html` for the root route (`/`)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Welcome to Student Enrollment System",
		})
	})

	// ------------------ Student Routes ------------------

	// Route for displaying students
	r.GET("/students", func(c *gin.Context) {
		students, err := controllers.FetchAllStudents(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching students: %v", err)
			return
		}
		c.HTML(http.StatusOK, "students.html", gin.H{"students": students})
	})

	// Route for adding a new student
	r.POST("/students/add", func(c *gin.Context) {
		// Extract form data
		name := c.PostForm("name")
		if name == "" {
			c.String(http.StatusBadRequest, "Name is required")
			return
		}

		age, err := strconv.Atoi(c.DefaultPostForm("age", "0"))
		if err != nil || age <= 0 {
			c.String(http.StatusBadRequest, "Invalid age provided")
			return
		}

		phone := c.PostForm("phone")
		if phone == "" {
			c.String(http.StatusBadRequest, "Phone number is required")
			return
		}

		image := c.PostForm("image")

		student := models.Student{
			Name:  name,
			Age:   age,
			Phone: phone,
			Image: image,
		}

		if err := controllers.AddStudent(db, student); err != nil {
			c.String(http.StatusInternalServerError, "Failed to add student: %v", err)
			return
		}

		c.Redirect(http.StatusSeeOther, "/students")
	})

	// Route for deleting a student
	r.POST("/students/delete", func(c *gin.Context) {
		studentID := c.PostForm("id")
		if studentID == "" {
			c.String(http.StatusBadRequest, "Student ID is required")
			return
		}

		// Convert string ID to integer
		intID, err := strconv.Atoi(studentID)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid student ID format")
			return
		}

		// Convert intID back to string before passing it to the controller
		if err := controllers.DeleteStudent(db, strconv.Itoa(intID)); err != nil {
			c.String(http.StatusInternalServerError, "Error deleting student: %v", err)
			return
		}
		c.Redirect(http.StatusSeeOther, "/students")
	})

	// Route for displaying the edit form for a specific student
	r.GET("/students/edit/:id", func(c *gin.Context) {
		studentID := c.Param("id")

		// Convert string ID to integer
		intID, err := strconv.Atoi(studentID)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid student ID format")
			return
		}

		student, err := controllers.FetchStudentByID(db, strconv.Itoa(intID))
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching student: %v", err)
			return
		}

		// Render the edit form with the student's current details
		c.HTML(http.StatusOK, "editStudent.html", gin.H{"student": student})
	})

	// Route for updating student details
	r.POST("/students/update", func(c *gin.Context) {
		// Get form values
		id := c.PostForm("id")
		name := c.PostForm("name")
		age, err := strconv.Atoi(c.DefaultPostForm("age", "0"))
		if err != nil || age <= 0 {
			c.String(http.StatusBadRequest, "Invalid age provided")
			return
		}
		phone := c.PostForm("phone")
		image := c.PostForm("image")

		// Validate inputs
		if name == "" || phone == "" {
			c.String(http.StatusBadRequest, "Name and Phone are required")
			return
		}

		// Convert string ID to integer
		intID, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid student ID format")
			return
		}

		// Create student model
		student := models.Student{
			ID:    intID,
			Name:  name,
			Age:   age,
			Phone: phone,
			Image: image,
		}

		// Capture both return values (error and result)
		if _, err := controllers.ModifyStudent(db, student); err != nil {
			c.String(http.StatusInternalServerError, "Failed to update student: %v", err)
			return
		}

		// Redirect back to the student list page
		c.Redirect(http.StatusSeeOther, "/students")
	})

	// ------------------ Course Routes ------------------

	// Route for listing all courses
	r.GET("/courses", func(c *gin.Context) {
		courses, err := controllers.FetchAllCourses(db) // Correct function for fetching courses
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching courses: %v", err)
			return
		}
		c.HTML(http.StatusOK, "courses.html", gin.H{"courses": courses})
	})

	// Route for adding a new course
	r.POST("/courses/add", func(c *gin.Context) {
		// Fetching form data
		name := c.PostForm("name")
		description := c.PostForm("description")
		creditsStr := c.DefaultPostForm("credit", "0") // Use 'credit' here to match the input field name in the form
	
		// Convert credit to integer
		credits, err := strconv.Atoi(creditsStr)
		if err != nil || credits <= 0 {
			c.String(http.StatusBadRequest, "Invalid credits value")
			return
		}
	
		// Add course
		err = controllers.AddCourse(db, name, description, credits)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to add course: %v", err)
			return
		}
	
		// Redirect to course list page
		c.Redirect(http.StatusSeeOther, "/courses")
	})
	// Route for displaying the edit form for a specific student
	r.GET("/courses/edit/:id", func(c *gin.Context) {
		courseID := c.Param("id")

		// Convert string ID to integer
		intID, err := strconv.Atoi(courseID)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid course ID format")
			return
		}

		course, err := controllers.FetchCourseByID(db, strconv.Itoa(intID))
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching course: %v", err)
			return
		}

		// Render the edit form with the student's current details
		c.HTML(http.StatusOK, "editCourse.html", gin.H{"course": course})
	})

	// Route to handle updating course details
r.POST("/courses/update", func(c *gin.Context) {
    id := c.PostForm("id")
    name := c.PostForm("name")
    description := c.PostForm("description")
    creditsStr := c.PostForm("credit")

    // Convert credits to integer
    credits, err := strconv.Atoi(creditsStr)
    if err != nil || credits <= 0 {
        c.String(http.StatusBadRequest, "Invalid credits value")
        return
    }

    // Update course in the database
    rowsAffected, err := controllers.UpdateCourse(db, id, name, description, credits)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to update course: %v", err)
        return
    }

    if rowsAffected == 0 {
        c.String(http.StatusNotFound, "Course not found")
        return
    }

    // Redirect to course list
    c.Redirect(http.StatusSeeOther, "/courses")
})

// Route to handle deleting a course
r.POST("/courses/delete", func(c *gin.Context) {
    id := c.PostForm("id")

    // Delete course from the database
    err := controllers.DeleteCourse(db, id)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to delete course: %v", err)
        return
    }

    // Redirect to course list
    c.Redirect(http.StatusSeeOther, "/courses")
})
// Enrollment routes
r.GET("/enrollments", func(c *gin.Context) { handlers.ShowEnrollmentsPage(c, db) })
r.POST("/enrollments/assign", func(c *gin.Context) { handlers.AssignEnrollment(c, db) })


	// Start the server on port 8080
	r.Run(":8080")
}
