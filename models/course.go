package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
	"fmt"
)

type Course struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     int    `json:"credits"`
}

// FetchAllCourses retrieves all courses from the database.
func FetchAllCourses(db *sql.DB) ([]Course, error) {
	rows, err := db.Query("SELECT id, name, description, credits FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Credits); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, rows.Err()
}

// AddCourse adds a new course to the database.
func AddCourse(db *sql.DB, name, description string, credits int) error {
	_, err := db.Exec("INSERT INTO courses (name, description, credits) VALUES (?, ?, ?)", name, description, credits)
	return err
}

func FindCourseByID(db *sql.DB, id int) (Course, error) {
    var course Course
    row := db.QueryRow("SELECT id, name, description, credits FROM courses WHERE id = ?", id)
    if err := row.Scan(&course.ID, &course.Name, &course.Description, &course.Credits); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return course, fmt.Errorf("course with ID %d not found", id)
        }
        return course, fmt.Errorf("error fetching course with ID %d: %w", id, err)
    }
    return course, nil
}


// UpdateCourse updates an existing course in the database.
func UpdateCourse(db *sql.DB, id int, name, description string, credits int) error {
	_, err := db.Exec("UPDATE courses SET name = ?, description = ?, credits = ? WHERE id = ?", name, description, credits, id)
	return err
}

// DeleteCourse deletes a course by its ID.
func DeleteCourse(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM courses WHERE id = ?", id)
	return err
}
