package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"sync"
)

func SendEmail(subject string, body string, to string, wg *sync.WaitGroup) bool {
	defer wg.Done()

	message := []byte("From: Contact <" + smtpConfig.Sender + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		body + "\r\n")

	// Create authentication credentials
	auth := smtp.PlainAuth("", smtpConfig.Sender, smtpConfig.Password, smtpConfig.SMTPHost)

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpConfig.SMTPHost,
	}

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpConfig.SMTPHost, smtpConfig.SMTPPort), tlsConfig)
	if err != nil {
		log.Printf("Failed to connect to the SMTP server: %v", err)
		return false
	}

	// Create the SMTP client
	client, err := smtp.NewClient(conn, smtpConfig.SMTPHost)
	if err != nil {
		log.Printf("Failed to create SMTP client: %v", err)
		return false
	}

	// Authenticate with the SMTP server
	if err := client.Auth(auth); err != nil {
		log.Printf("SMTP authentication failed: %v", err)
		return false
	}

	// Set the sender and recipient
	if err := client.Mail(smtpConfig.Sender); err != nil {
		log.Printf("Failed to set sender: %v", err)
		return false
	}
	if err := client.Rcpt(to); err != nil {
		log.Printf("Failed to set recipient: %v", err)
		return false
	}

	// Send the email message
	w, err := client.Data()
	if err != nil {
		log.Printf("Failed to open data writer: %v", err)
		return false
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.Write(message)

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Failed to write email message: %v", err)
		return false
	}

	err = w.Close()
	if err != nil {
		log.Printf("Failed to close data writer: %v", err)
		return false
	}

	// Close the connection to the SMTP server
	client.Quit()

	log.Println("Email sent successfully!")
	return true
}
