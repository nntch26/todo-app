// Package domain holds the core data types. It must not import Fiber or GORM.
package domain

import "time"

type User struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
}
