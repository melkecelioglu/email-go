package main

type RequestBody struct {
	RecipientEmail string `json:"email"`
	PhoneNumber    string `json:"phone"`
	Name           string `json:"name"`
}

type SMTPConfig struct {
	SMTPHost string `yaml:"smtpHost"`
	SMTPPort int    `yaml:"smtpPort"`
	Sender   string `yaml:"smtpEmail"`
	Password string `yaml:"smtpPassword"`
}
