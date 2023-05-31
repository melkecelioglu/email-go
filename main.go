package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

// Global SMTP server configuration variable
var smtpConfig SMTPConfig

// Load SMTP server configuration from the YAML file
func loadSMTPConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &smtpConfig)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	err := loadSMTPConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load SMTP configuration: %v", err)
	}

	// Define the endpoint for sending the email
	http.HandleFunc("/api/book", SendEmailHandler)

	// Start the HTTP server
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
