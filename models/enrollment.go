package models

import "database/sql"

// Enrollment represents an enrollment record.
type Enrollment struct {
    StudentID   int    
    CourseID    int    
    StudentName string 
    CourseName  string 
    EnrolledAt  string 
}

// GetEnrollments retrieves all enrollments from the database.
func GetEnrollments(db *sql.DB) ([]Enrollment, error) {
    rows, err := db.Query(`
        SELECT e.student_id, e.course_id, s.name AS student_name, c.name AS course_name, e.enrolled_at
        FROM enrollments e
        JOIN students s ON e.student_id = s.id
        JOIN courses c ON e.course_id = c.id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var enrollments []Enrollment
    for rows.Next() {
        var enrollment Enrollment
        err := rows.Scan(&enrollment.StudentID, &enrollment.CourseID, &enrollment.StudentName, &enrollment.CourseName, &enrollment.EnrolledAt)
        if err != nil {
            return nil, err
        }
        enrollments = append(enrollments, enrollment)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return enrollments, nil
}

// CreateEnrollment assigns a student to a course.
func CreateEnrollment(db *sql.DB, studentID, courseID int) error {
    _, err := db.Exec("INSERT INTO enrollments (student_id, course_id, enrolled_at) VALUES (?, ?, NOW())", studentID, courseID)
    return err
}

