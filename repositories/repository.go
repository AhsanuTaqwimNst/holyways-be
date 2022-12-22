package repositories

import "gorm.io/gorm"

// Declare repository struct is here

type repository struct {
	db *gorm.DB
}
