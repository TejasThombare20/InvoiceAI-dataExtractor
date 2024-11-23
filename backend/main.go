package main

import (
	"log"

	"github.com/TejasThombare20/backend/config"
	"github.com/TejasThombare20/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MongoDB connection
	if err := config.ConnectDB(); err != nil {
		// return fmt.Errorf("failed to connect to MongoDB: %v", err)
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	log.Println("MongoDB connected successfully")

	// Initialize Gemini client
	if err := config.InitGemini(); err != nil {
		// return fmt.Errorf("failed to initialize Gemini client: %v", err)
		log.Fatalf("failed to initialize Gemini client: %v", err)
	}
	log.Println("Gemini client initialized successfully")

	// Create Gin router
	router := gin.Default()

	// router.Use(gin.Logger())

	// Setup CORS
	router.MaxMultipartMemory = 8 << 20
	router.Use(cors.Default())

	// Initialize routes
	routes.SetupRoutes(router)

	// Start server
	router.Run(":8080")
}
