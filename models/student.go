package models

import (
	"database/sql"
	"fmt"
)

// Student represents a student in the system
type Student struct {
	ID    int    // Changed from id to ID for exportability
	Name  string
	Age   int
	Phone string
	Image string
}

// GetAllStudents fetches all students from the database
func GetAllStudents(db *sql.DB) ([]Student, error) {
	rows, err := db.Query("SELECT id, name, age, phone, image FROM students")
	if err != nil {
		return nil, fmt.Errorf("GetAllStudents failed: %v", err)
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Phone, &student.Image)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student row: %v", err)
		}
		students = append(students, student)
	}
	return students, rows.Err()
}

// AddStudent adds a new student to the database
func AddStudent(db *sql.DB, student Student) error {
	_, err := db.Exec("INSERT INTO students (name, age, phone, image) VALUES (?, ?, ?, ?)", 
		student.Name, student.Age, student.Phone, student.Image)
	if err != nil {
		return fmt.Errorf("AddStudent failed: %v", err)
	}
	return nil
}

// GetStudentByID retrieves a student by ID from the database
func GetStudentByID(db *sql.DB, id int) (*Student, error) {
	var student Student
	err := db.QueryRow("SELECT id, name, age, phone, image FROM students WHERE id = ?", id).
		Scan(&student.ID, &student.Name, &student.Age, &student.Phone, &student.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve student by ID: %v", err)
	}
	return &student, nil
}

// UpdateStudent updates a student's information in the database
func UpdateStudent(db *sql.DB, student Student) (int64, error) {
	result, err := db.Exec("UPDATE students SET name = ?, age = ?, phone = ?, image = ? WHERE id = ?", 
		student.Name, student.Age, student.Phone, student.Image, student.ID)
	if err != nil {
		return 0, fmt.Errorf("UpdateStudent failed: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	return rowsAffected, nil
}

// DeleteStudent deletes a student from the database by ID
func DeleteStudent(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteStudent failed for ID %d: %v", id, err)
	}
	return nil
}
