package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"piszcodziennie/internal/database"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
)

type Config struct {
	db             *database.Queries
	port           string
	firebaseClient *auth.Client
	env            string
	SMTPEnvs       SMTPEnvs
}

type SMTPEnvs struct {
	Email       string
	SMTP_PWD    string
	SMTP_SERVER string
	SMTP_PORT   string
}

func getSMTPEnvs() (SMTPEnvs, error) {
	email := os.Getenv("SMTP_EMAIL")
	smtp_pwd := os.Getenv("SMTP_PWD")
	smtp_server := os.Getenv("SMTP_SERVER")
	smtp_port := os.Getenv("SMTP_PORT")
	if email == "" || smtp_pwd == "" || smtp_server == "" || smtp_port == "" {
		return SMTPEnvs{}, fmt.Errorf("SMTP_EMAIL, SMTP_PWD, SMTP_SERVER and SMTP_PORT must be set")
	}
	return SMTPEnvs{
		Email:       email,
		SMTP_PWD:    smtp_pwd,
		SMTP_SERVER: smtp_server,
		SMTP_PORT:   smtp_port,
	}, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	serviceAccountKey := "piszcodziennie-firebase.json"

	opt := option.WithCredentialsFile(serviceAccountKey)
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Please set PORT environment variable")
	}

	smtpEnvs, err := getSMTPEnvs()
	if err != nil {
		log.Fatal(err)
	}

	cfg := Config{
		db:             dbQueries,
		firebaseClient: client,
		port:           PORT,
		env:            env,
		SMTPEnvs:       smtpEnvs,
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/health", cfg.healthHandler)
	r.Get("/", cfg.home)
	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "register_form.html", nil)
	})
	r.Post("/register", cfg.registerHandler)
	r.Post("/login", cfg.loginHandler)
	r.Get("/protected", cfg.authMiddleware(cfg.protectedHandler))

	r.Post("/reset", cfg.reset)

	srv := &http.Server{
		Addr:              ":" + cfg.port,
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
	}
	log.Println("Server is running on port " + cfg.port)
	log.Fatal(srv.ListenAndServe())
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	fmt.Println("templates/" + tmpl)
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (cfg *Config) home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		renderTemplate(w, "login_form.html", nil)
		return
	}

	// Verify token from cookie
	token, err := cfg.firebaseClient.VerifyIDToken(context.Background(), cookie.Value)
	if err != nil {
		renderTemplate(w, "login_form.html", nil)
		return
	}

	renderTemplate(w, "dashboard.html", map[string]string{"email": token.Claims["email"].(string)})
}
