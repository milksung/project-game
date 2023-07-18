package repository

import (
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}
