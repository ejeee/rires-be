package main

import (
	"fmt"
	"rires-be/config"
	"rires-be/internal/models"
	"rires-be/pkg/database"
)

func main() {
	// Load config
	if err := config.LoadConfig(); err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	// Connect to database
	if err := database.Connect(config.AppConfig.GetDSN()); err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer database.CloseDB()

	// Test query
	var menus []models.Menu
	if err := database.DB.Where("hapus = ?", 0).Order("parent_id ASC, urutan ASC").Limit(5).Find(&menus).Error; err != nil {
		fmt.Println("Query error:", err)
		return
	}

	fmt.Printf("Found %d menus:\n", len(menus))
	for _, menu := range menus {
		fmt.Printf("- ID: %d, Name: %s, URL: %s\n", menu.ID, menu.NamaMenu, menu.URLMenu)
	}
}
