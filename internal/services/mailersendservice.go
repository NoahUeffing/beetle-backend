package services

import (
	"beetle/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MailerSend API URL
const mailerSendURL = "https://api.mailersend.com/v1/email"

type MailerSendService struct {
	Config config.EmailConfig
}

func NewMailerSendService(config config.EmailConfig) *MailerSendService {
	return &MailerSendService{
		Config: config,
	}
}

// EmailRequest defines the payload structure
type EmailRequest struct {
	From    EmailAddress   `json:"from"`
	To      []EmailAddress `json:"to"`
	Subject string         `json:"subject"`
	Text    string         `json:"text"`
	HTML    string         `json:"html"`
}

// EmailAddress holds an email and optional name
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// sendEmail sends an email via MailerSend API
func (ms *MailerSendService) Send(fromEmail, toEmail, subject, textBody, htmlBody string) error {
	email := EmailRequest{
		From: EmailAddress{
			Email: fromEmail,
			Name:  "Beetle",
		},
		To: []EmailAddress{
			{
				Email: toEmail,
				//Name:  "", // Optional
			},
		},
		Subject: subject,
		Text:    textBody,
		HTML:    htmlBody,
	}

	payload, err := json.Marshal(email)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", mailerSendURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", "Bearer "+ms.Config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Email sent successfully.")
		return nil
	}

	// Print the error response from MailerSend
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("MailerSend error response:", string(body))

	return fmt.Errorf("failed to send email: %s", resp.Status)
}
