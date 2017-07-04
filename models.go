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
	Address  string
	Timezone string
}

// Semester
type Semester struct {
	gorm.Model
	Slug string `gorm:"unique_index"`
	Name string
	Year string
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
	Campus Campus
	Room   string
	Sun    bool
	Mon    bool
	Tues   bool
	Wed    bool
	Thurs  bool
	Fri    bool
	Sat    bool
	Start  time.Time
	End    time.Time
}

// ClassType
type ClassType struct {
	gorm.Model
	Name string
}

// Section
type Section struct {
	gorm.Model
	Course    Course
	Sessions  []Session
	Professor string
	ClassType ClassType
}
