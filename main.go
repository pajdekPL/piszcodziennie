package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	supa "github.com/supabase-community/supabase-go"
)

type Config struct {
	port           string
	supabaseClient *supa.Client
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("Please set SUPABASE_URL and SUPABASE_KEY environment variables")
	}
	client, err := supa.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatal("Failed to initialize Supabase client:", err)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Please set PORT environment variable")
	}

	cfg := Config{
		supabaseClient: client,
	}

	r := gin.Default()

	r.POST("/register", cfg.handlerRegisterUser)
	r.POST("/login", cfg.handlerLoginUser)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	err = r.Run(fmt.Sprintf(":%s", cfg.port))
	if err != nil {
		log.Fatal(err)
	}

}
