package main

import (
	"backend/api"
	"backend/internal/dbrepo"
	"backend/pkg/auth"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"github.com/joho/godotenv"
)

const port = 8080

func main() {
	// set application config
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file. Using default values")
	}
	var app api.Application

	// read from command line
	flag.StringVar(&app.DSN, "dsn", getEnv("DSN", "host=localhost port=5432 user=postgres password=postgres dbname=myserver sslmode=disable timezone=UTC connect_timeout=5"), "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", getEnv("JWT_SECRET", "verysecret"), "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", getEnv("JWT_ISSUER", "example.com"), "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", getEnv("JWT_AUDIENCE", "example.com"), "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", getEnv("COOKIE_DOMAIN", "localhost"), "cookie domain")
	flag.StringVar(&app.Domain, "domain", getEnv("DOMAIN", "example.com"), "domain")

	flag.Parse()
	// Initialize the database connection
	repo := &dbrepo.PostgresDBRepo{}
	app.DB = repo

	db, err := app.DB.ConnectToDB(app.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}