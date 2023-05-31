package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Handler function for sending email
func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request /api/contact")

	// Parse the JSON request body
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Send the email concurrently
	var wg sync.WaitGroup
	wg.Add(1)

	// Compose the email message
	subject := fmt.Sprintf("Contact received for %s", requestBody.Name)
	body, err := getEmailBody("email-template.html", requestBody)
	if err != nil {
		log.Println("Failed to read email template:", err)
		http.Error(w, "Failed to read email template", http.StatusInternalServerError)
		return
	}

	go SendEmail(subject, body, requestBody.RecipientEmail, &wg)

	// Send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

func getEmailBody(templateFile string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
