package main

import (
	"backend/api"
	"backend/internal/dbrepo"
	"backend/pkg/auth"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx  driver for database/sql
	"github.com/joho/godotenv"
)

var port int

func main() {
	// set application config
	var app api.Application
	var err error

	// Locate the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Construct the path to the .env file
	envPath := filepath.Join(dir, ".env")

	// Check if the .env file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Println(".env file not found, using default values")
	} else {
		log.Println(".env file found, loading environment variables")
		// Load environment variables from .env file if present
		err = godotenv.Load(envPath)
		if err != nil {
			log.Fatal("Error loading .env file. It is not readable.")
		}
	}
	// Load configuration from environment variables
	app.DSN = os.Getenv("DSN")
	app.JWTSecret = os.Getenv("JWT_SECRET")

	// Set default values for the application configuration
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.IntVar(&port, "port", 8080, "API server port")

	flag.Parse()
	// Initialize the database connection
	if app.DSN == "" {
		log.Fatal("DSN environment variable is not set")
	}
	if app.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	db, err := sql.Open("pgx", app.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Test the database connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Connected to the database")

	// Set up the database connection pool
	defer db.Close()

	// Assign the database repo to app.DB
	app.DB = &dbrepo.PostgresDBRepo{DB: db}

	app.Auth = auth.Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	// Start a web server
	fmt.Printf("Starting server on port %d\n", port)

	// Set up your application's router/handler
	handler := app.Routes()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))

}
