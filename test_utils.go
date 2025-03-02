package main

// import (
// 	"log"

// 	"github.com/gin-gonic/gin"
// 	supa "github.com/supabase-community/supabase-go"
// )

// func CreateTestConfig() *Config {
// 	supabaseURL := "http://localhost:54321"
// 	supabaseKey := "your-local-anon-key"

// 	client, err := supa.NewClient(supabaseURL, supabaseKey, nil)

// 	if err != nil {
// 		log.Fatal("Failed to initialize Supabase client:", err)
// 	}

// 	return &Config{
// 		supabaseClient: client,
// 		port:           "8080",
// 	}
// }
// func SetupTestRouter(cfg *Config) *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()
// 	router.POST("/register", cfg.handlerRegisterUser)
// 	return router
// }
