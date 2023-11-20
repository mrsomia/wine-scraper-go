package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string
	ProductURL []ProductURL
	Price      []Price
}

type ProductURL struct {
	gorm.Model
	ProductID uint
	Store   string
	URL     string
}

type Price struct {
	gorm.Model
	ProductID uint
  Store string
	Price   float64
}

func GetDBConnection(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{}, &ProductURL{}, &Price{})
	return db
}
