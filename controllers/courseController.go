package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"student-enrollment-system/models"
)

// FetchAllCourses fetches all courses from the database.
func FetchAllCourses(db *sql.DB) ([]models.Course, error) {
	rows, err := db.Query("SELECT id, name, description, credits FROM courses")
	if err != nil {
		log.Println("Error fetching courses:", err)
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Credits); err != nil {
			log.Println("Error scanning course:", err)
			return nil, err
		}
		courses = append(courses, course)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return courses, nil
}
// AddCourse adds a new course to the database.
func AddCourse(db *sql.DB, name, description string, credits int) error {
	_, err := db.Exec("INSERT INTO courses (name, description, credits) VALUES (?, ?, ?)", name, description, credits)
	if err != nil {
		log.Println("Error adding course:", err)
		return err
	}
	return nil
}

// FetchCourseByID retrieves a course by ID from the database
// FetchCourseByID retrieves a course by ID from the database
func FetchCourseByID(db *sql.DB, id string) (*models.Course, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	// Corrected SQL query with proper column names
	row := db.QueryRow("SELECT id, name, description, credits FROM courses WHERE id = ?", intID)

	var course models.Course
	err = row.Scan(&course.ID, &course.Name, &course.Description, &course.Credits)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Course with ID %d not found", intID)
			return nil, fmt.Errorf("course with ID %d not found", intID)
		}
		log.Printf("Error fetching course by ID: %v", err)
		return nil, fmt.Errorf("could not fetch course with ID %d: %v", intID, err)
	}

	return &course, nil
}


// UpdateCourse updates an existing course's details in the database
// UpdateCourse updates an existing course's details in the database
func UpdateCourse(db *sql.DB, id, name, description string, credits int) (int64, error) {
	// Corrected SQL query with proper column names
	result, err := db.Exec("UPDATE courses SET name = ?, description = ?, credits = ? WHERE id = ?",
		name, description, credits, id)
	if err != nil {
		log.Printf("Error updating course: %v", err)
		return 0, fmt.Errorf("could not update course: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows: %v", err)
		return 0, fmt.Errorf("could not get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		log.Printf("No rows affected when updating course with ID: %s", id)
		return 0, fmt.Errorf("no rows affected when updating course with ID: %s", id)
	}

	log.Printf("Successfully updated course with ID: %s", id)
	return rowsAffected, nil
}


// DeleteCourse deletes a course from the database by ID
// DeleteCourse deletes a course from the database by ID
func DeleteCourse(db *sql.DB, id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return fmt.Errorf("invalid ID format: %v", err)
	}

	// Corrected SQL query with proper column names
	_, err = db.Exec("DELETE FROM courses WHERE id = ?", intID)
	if err != nil {
		log.Printf("Error deleting course with ID %d: %v", intID, err)
		return fmt.Errorf("could not delete course with ID %d: %v", intID, err)
	}

	log.Printf("Successfully deleted course with ID: %d", intID)
	return nil
}

