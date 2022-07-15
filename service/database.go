package service

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	// shared database connection
	DB *gorm.DB
)

// Opens the database connection using ConnectionString. Database connection
// can be accessed using globally accessible DB variable.
func Connect(connectionString string, tablePrefix string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tablePrefix,
		},
	})
	return err
}

// Auto migrates models for database
func Migrate() error {
	return DB.AutoMigrate(&MatchingTable{})
}
