package courseModel

import "time"

type Course struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DurationInHours int       `json:"durationInHours"`
	StartDate       time.Time `json:"startDate"`
	Instructor      string    `json:"instructor"`
}
