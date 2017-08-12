package schedules

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Campus
type Campus struct {
	gorm.Model
	Slug     string `gorm:"unique_index"`
	Name     string
	Timezone string
}

// Semester
type Semester struct {
	//gorm.Model
	//Slug string `gorm:"unique_index"`
	Name string
	//Year string
}

// Course
type Course struct {
	gorm.Model
	Slug string `gorm:"unique_index"`
	Name string
}

// Session
type Session struct {
	gorm.Model
	Campus       Campus
	Location     string
	SectionRefer uint
	Sun          bool `gorm:"default:0"`
	Mon          bool `gorm:"default:0"`
	Tues         bool `gorm:"default:0"`
	Wed          bool `gorm:"default:0"`
	Thurs        bool `gorm:"default:0"`
	Fri          bool `gorm:"default:0"`
	Sat          bool `gorm:"default:0"`
	Start        time.Time
	End          time.Time
}

// Section
type Section struct {
	gorm.Model
	Semester     string    `gorm:"not null"` // Semester the section was taught in
	Course       string    `gorm:"not null"` // Course the section belongs to
	Section      string    `gorm:"not null"` // Section ID within the course (ex: "MGMT 313 001" --> "001")
	CRN          int       `gorm:"not null"` // Course Registration Number
	LinkCode     string    // ?
	ScheduleType string    `gorm:"not null"` // What kind of class this represents (ex: Lecture, Internship, Recitation)
	SectionTitle string    `gorm:"not null"` // Longform Title of the Course (ex: Organizational Behavior)
	Instructor   string    `gorm:"not null"` // Name of the Instructor teaching the course
	StartDate    string    `gorm:"not null"` // Start date of the course
	EndDate      string    `gorm:"not null"` // End date of the course
	Days         string    // Days the section meets
	Times        string    // Times on those days the section meets
	Sessions     []Session `gorm:"ForeignKey:SectionRefer"` // Sessions the section meets at (time/places/etc)
	SecDays      string    // Secondary Days
	SecTimes     string    // Secondary Times
	TertDays     string    // Tertiary Days
	TertTimes    string    // Tertiary Times
	Location     string    `gorm:"not null"` // Where the section meets
	RoomCapacity int       // Seating capacity of the location
	Status       string    `gorm:"not null"`                                         // Open, Waitlisted, Cancelled, Restricted
	CampusName   string    `gorm:"not null"`                                         // Campus Name
	Campus       Campus    `gorm:"ForeignKey:CampusName;AssociationForeignKey:Name"` // What campus the location is on
	XList        bool      // ?
	XListNumber  string    // ?
	SectionNotes string    // Notes from the registrar's office
	MinCRHours   int       // Minimum number of credits that can be applied to the course
	Limit        int       // Seating limit for the course
	Enrolled     int       // Number of enrolled students
	SeatsAvail   int       // Number of available seats
	Waitlist     int       // Number of students on the waitlist
}
