package domain

import "time"

type Todo struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string     `gorm:"type:uuid;not null;index"`
	Title       string     `gorm:"type:varchar(200);not null"`
	Description string     `gorm:"type:varchar(2000);not null;default:''"`
	DueDate     *time.Time `gorm:"type:date"` // nil = no due date
	Done        bool       `gorm:"not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User User `gorm:"constraint:OnDelete:CASCADE"`
}

// TodoPatch is a partial update; nil field = leave unchanged
type TodoPatch struct {
	Title       *string
	Description *string
	Done        *bool
	DueDate     *DueDateChange
}

// DueDateChange: Date nil = remove the due date
type DueDateChange struct {
	Date *time.Time
}
