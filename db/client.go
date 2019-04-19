package db

import (
	"github.com/jinzhu/gorm"
)

type Dao struct {
	DB *gorm.DB
}
