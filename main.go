package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Phone string `json:"phone"`
	Image string `json:"image"`
}

type Course struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     int    `json:"credits"`
}


func initDB() {
	var err error

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/student_enrollment_system")
	if err != nil {
		panic(err)
	}
}

func main() {
	initDB()

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", homePage)
	r.GET("/students", listStudents)
	r.POST("/students/add", addStudent)
	r.GET("/students/edit/:id", getStudentForEdit) 
	r.POST("/students/update", updateStudent)   
	r.POST("/students/delete", deleteStudent)

	r.GET("/courses", listCourses)
	r.GET("/courses/edit/:id", getCourseForEdit)
	r.POST("/courses/update", updateCourse)     
	r.POST("/courses/add", addCourse)
	r.POST("/courses/delete", deleteCourse)

	r.GET("/enrollments", listEnrollments)
	r.POST("/enrollments/assign", assignEnrollment)

	r.Run(":8080")
}
func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// Students Handlers
func listStudents(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, age, phone, image FROM students")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var students []map[string]interface{}
	for rows.Next() {
		var id int
		var name, phone, image string
		var age int
		err := rows.Scan(&id, &name, &age, &phone, &image)
		if err != nil {
			panic(err)
		}

		students = append(students, map[string]interface{}{
			"id":    id,
			"name":  name,
			"age":   age,
			"phone": phone,
			"image": image,
		})
	}

	if err := rows.Err(); err != nil {
		c.String(http.StatusInternalServerError, "Error during rows iteration: %v", err)
		return
	}

	c.HTML(http.StatusOK, "students.html", gin.H{"students": students})
}

func addStudent(c *gin.Context) {
	name := c.PostForm("name")
	age, _ := strconv.Atoi(c.PostForm("age"))
	phone := c.PostForm("phone")
	image := c.PostForm("image") 

	_, _ = db.Exec("INSERT INTO students (name, age, phone, image) VALUES (?, ?, ?, ?)", name, age, phone, image)
	c.Redirect(http.StatusSeeOther, "/students")
}

func getStudentForEdit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}

	var student Student
	err := db.QueryRow("SELECT id, name, age, phone, image FROM students WHERE id = ?", id).
		Scan(&student.ID, &student.Name, &student.Age, &student.Phone, &student.Image)

	if err != nil {
		fmt.Println("Error fetching student: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	fmt.Println(student)

	c.HTML(http.StatusOK, "editStudent.html", gin.H{
		"student": student,
	})
}

func updateStudent(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	age, err := strconv.Atoi(c.PostForm("age"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age format"})
		return
	}
	phone := c.PostForm("phone")
	image := c.PostForm("image")
	result, err := db.Exec("UPDATE students SET name = ?, age = ?, phone = ?, image = ? WHERE id = ?", name, age, phone, image, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found or not updated"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/students")
}

func deleteStudent(c *gin.Context) {
	id := c.PostForm("id")
	_, _ = db.Exec("DELETE FROM students WHERE id = ?", id)
	c.Redirect(http.StatusSeeOther, "/students")
}

// Courses Handlers
func listCourses(c *gin.Context) {

	rows, err := db.Query("SELECT id, name, description, credits FROM courses")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve courses: %v", err)
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Credits); err != nil {
			c.String(http.StatusInternalServerError, "Failed to scan course: %v", err)
			return
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		c.String(http.StatusInternalServerError, "Error during rows iteration: %v", err)
		return
	}

	c.HTML(http.StatusOK, "courses.html", gin.H{"courses": courses})
}

func addCourse(c *gin.Context) {

	courseName := c.PostForm("course_name") 
	courseDescription := c.PostForm("course_description") 
	courseCredits, err := strconv.Atoi(c.PostForm("course_credit")) 
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid course credit value: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO courses (name, description, credits) VALUES (?, ?, ?)",
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

	var course Course
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

// Enrollments Handlers
func listEnrollments(c *gin.Context) {

    enrollmentsRows, err := db.Query(`
        SELECT e.id, s.name AS student_name, c.name AS course_name, e.enrolled_at
        FROM enrollments e
        JOIN students s ON e.student_id = s.id
        JOIN courses c ON e.course_id = c.id
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query enrollments"})
        return
    }
    defer enrollmentsRows.Close()

    var enrollments []map[string]interface{}
    for enrollmentsRows.Next() {
        var id int
        var studentName, courseName, enrolledAt string

        err := enrollmentsRows.Scan(&id, &studentName, &courseName, &enrolledAt)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan enrollments"})
            return
        }

        enrollments = append(enrollments, map[string]interface{}{
            "id":           id,
            "student_name": studentName,
            "course_name":  courseName,
            "enrolled_at":  enrolledAt,
        })
    }

    if err := enrollmentsRows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating enrollments rows"})
        return
    }

    c.HTML(http.StatusOK, "enrollments.html", gin.H{"enrollments": enrollments})
}


// Enrollments Handlers
func assignEnrollment(c *gin.Context) {

	studentID := c.PostForm("student_id")
	courseID := c.PostForm("course_id")

	if studentID == "" || courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student ID and Course ID are required"})
		return
	}
	_, err := db.Exec("INSERT INTO enrollments (student_id, course_id, enrolled_at) VALUES (?, ?, NOW())", studentID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign enrollment"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/enrollments")
}