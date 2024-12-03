package controllers

import (
	"database/sql"
	"student-enrollment-system/models"
	"fmt"
	"log"
	"strconv"
)

// FetchAllStudents retrieves all students from the database
func FetchAllStudents(db *sql.DB) ([]models.Student, error) {
	rows, err := db.Query("SELECT id, name, age, phone, image FROM students")
	if err != nil {
		log.Printf("Error querying students: %v", err)
		return nil, fmt.Errorf("could not fetch students: %v", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		// Scan into exported fields (ID instead of id)
		err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Phone, &student.Image)
		if err != nil {
			log.Printf("Error scanning student row: %v", err)
			return nil, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, fmt.Errorf("error iterating over student rows: %v", err)
	}

	return students, nil
}

// AddStudent adds a new student to the database
func AddStudent(db *sql.DB, student models.Student) error {
	_, err := db.Exec("INSERT INTO students (name, age, phone, image) VALUES (?, ?, ?, ?)",
		student.Name, student.Age, student.Phone, student.Image)
	if err != nil {
		log.Printf("Error inserting student: %v", err)
		return fmt.Errorf("could not insert student: %v", err)
	}
	log.Printf("Student added successfully: %v", student.Name)
	return nil
}

// FetchStudentByID retrieves a student by ID from the database
func FetchStudentByID(db *sql.DB, id string) (*models.Student, error) {
	// Convert string ID to int
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	// Call the model function to fetch the student by ID
	student, err := models.GetStudentByID(db, intID)
	if err != nil {
		log.Printf("Error fetching student by ID: %v", err)
		return nil, fmt.Errorf("could not fetch student with ID %d: %v", intID, err)
	}

	return student, nil
}

// ModifyStudent updates an existing student's details in the database
func ModifyStudent(db *sql.DB, student models.Student) (int64, error) {
	// Use the model function to update the student
	rowsAffected, err := models.UpdateStudent(db, student)
	if err != nil {
		log.Printf("Error updating student: %v", err)
		return 0, fmt.Errorf("could not update student: %v", err)
	}
	if rowsAffected == 0 {
		log.Printf("No rows affected when updating student with ID: %d", student.ID)
		return 0, fmt.Errorf("no changes made to student with ID: %d", student.ID)
	}
	log.Printf("Successfully updated student with ID: %d", student.ID)
	return rowsAffected, nil
}

// DeleteStudent deletes a student from the database by ID
func DeleteStudent(db *sql.DB, id string) error {
	// Convert string ID to int
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return fmt.Errorf("invalid ID format: %v", err)
	}

	// Call the model function to delete the student by ID
	err = models.DeleteStudent(db, intID)
	if err != nil {
		log.Printf("Error deleting student with ID %d: %v", intID, err)
		return fmt.Errorf("could not delete student with ID %d: %v", intID, err)
	}

	log.Printf("Successfully deleted student with ID: %d", intID)
	return nil
}
