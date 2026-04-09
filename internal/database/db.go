package database

import (
	"github.com/yogibala/auto-apply/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// This creates a file named 'auto_apply.db' in your root folder
	DB, err = gorm.Open(sqlite.Open("data/auto_apply.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Logic: Automatically create the table based on our 'JobApplication' struct
	DB.AutoMigrate(&models.JobApplication{})
}
