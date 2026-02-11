package smtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func sendEmail(to, subject, body string) error {
	lambdaURL := os.Getenv("EMAIL_LAMBDA_URL") // Lambda Function URL

	payload := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", lambdaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-EMAIL-TOKEN", os.Getenv("EMAIL_LAMBDA_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("email lambda failed with status %d", resp.StatusCode)
	}

	return nil
}

func SendEmailOTP(to, otp string) error {
	subject := "Your OTP Code"

	// HTML-styled email body
	body := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Your OTP Code</title>
	</head>
	<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background-color: #f4f4f4; padding: 20px; border-radius: 8px;">
			<h2 style="color: #2c3e50; text-align: center;">Your OTP Code</h2>
			<p>Dear User,</p>
			<p>Your One-Time Password (OTP) is:</p>
			<p style="font-size: 24px; font-weight: bold; color: #e74c3c; text-align: center; background-color: #fff; padding: 10px; border: 1px solid #ddd; border-radius: 4px;">
				` + otp + `
			</p>
			<p>This OTP is valid for the next 5 minutes. Please do not share it with anyone.</p>
			<p>If you didn’t request this, please ignore this email.</p>
			<p style="text-align: center; margin-top: 20px;">
				<a href="https://yourwebsite.com" style="color: #3498db; text-decoration: none;">Visit our website</a>
			</p>
			<p>Best regards,<br>Sahil</p>
		</div>
	</body>
	</html>
	`

	err := sendEmail(to, subject, body)
	if err != nil {
		return err
	}
	fmt.Println("Email sent successfully to", to)
	return nil
}
