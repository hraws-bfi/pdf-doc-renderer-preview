package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	AllowedOrigins []string
	TemplatesDir   string
	ServerAddr     string
	DMS            DMSConfig
}

// DMSConfig holds DMS-specific configuration
type DMSConfig struct {
	APIURL    string
	APISecret string
}

var config Config

func init() {
	loadConfig()
}

func loadConfig() {
	// Load .env file (looks in current dir and parent dir)
	if err := godotenv.Load(); err != nil {
		// Try parent directory (for when running from render-api folder)
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("no .env file found, using environment variables")
		}
	}

	// Set defaults
	config.TemplatesDir = "../templates"
	config.ServerAddr = ":8080"

	// Load allowed origins
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		config.AllowedOrigins = []string{}
	} else {
		for _, o := range strings.Split(origins, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				config.AllowedOrigins = append(config.AllowedOrigins, o)
			}
		}
	}
	log.Printf("allowed CORS origins: %v", config.AllowedOrigins)

	// Load DMS configuration
	config.DMS.APIURL = os.Getenv("DMS_API_URL")
	config.DMS.APISecret = os.Getenv("DMS_API_SECRET")
	if config.DMS.APIURL != "" {
		log.Printf("DMS API configured: %s", config.DMS.APIURL)
	} else {
		log.Println("DMS API not configured (DMS_API_URL not set)")
	}
}
