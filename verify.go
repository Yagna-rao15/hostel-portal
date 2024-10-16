package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	// "golang.org/x/crypto/bcrypt"
	"log"
	"math/big"
	"net/smtp"
	"time"
)

var otpStore = map[string]string{}     // In-memory store for OTPs (email -> OTP)
var otpExpiry = map[string]time.Time{} // In-memory store for OTP expiry

func generateOTP() string {
	const otpLength = 6
	const charset = "0123456789"
	otp := ""
	for i := 0; i < otpLength; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		otp += string(charset[randomIndex.Int64()])
	}
	return otp
}

func sendEmail(email string, otp string) error {
	from := "theshamelescreature@gmail.com" // Sender's email
	password := "ynkzvqhwgtsdjmay"          // Sender's email password or app-specific password
	to := email                             // Recipient's email
	smtpHost := "smtp.gmail.com"            // SMTP server (Gmail, Outlook, etc.)
	smtpPort := "587"                       // SMTP port

	message := []byte(fmt.Sprintf("Subject: Your OTP Code\n\nYour OTP code is: %s", otp))
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}

func sendOTP(email string) error {
	otp := generateOTP()
	otpStore[email] = otp                              // Store OTP in memory
	otpExpiry[email] = time.Now().Add(5 * time.Minute) // Set expiry time (5 minutes)
	log.Println("Preparing to send email...")
	err := sendEmail(email, otp) // Send OTP via email
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	fmt.Println("OTP sent to:", email)
	return nil
}

func verifyOTP(email string, userOTP string) bool {
	storedOTP, exists := otpStore[email]
	if !exists {
		fmt.Println("No OTP found for this email.")
		return false
	}

	if time.Now().After(otpExpiry[email]) {
		fmt.Println("OTP has expired.")
		delete(otpStore, email)
		delete(otpExpiry, email)
		return false
	}

	if storedOTP == userOTP {
		fmt.Println("OTP verified successfully!")
		delete(otpStore, email)
		delete(otpExpiry, email)
		return true
	}
	fmt.Println("Incorrect OTP.")
	return false
}

func helper_main() {
	email := "yagnarao15@gmail.com"

	err := sendOTP(email)
	if err != nil {
		log.Fatal("Error sending OTP:", err)
	}

	var userOTP string
	fmt.Println("Enter the OTP received on your email:")
	fmt.Scanln(&userOTP)

	if verifyOTP(email, userOTP) {
		fmt.Println("User authenticated successfully.")
	} else {
		fmt.Println("OTP verification failed.")
	}
}

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
// 	return string(bytes), err
// }

func sendEmail1(email string, otp string) error {
	from := "theshamelescreature@gmail.com" // Sender's email
	password := "Yagnarao!15"               // Sender's email password or app-specific password
	to := email                             // Recipient's email
	smtpHost := "smtp.gmail.com"            // SMTP server (Gmail, Outlook, etc.)
	smtpPort := "465"                       // SMTP port

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Connect to the SMTP server with TLS
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to the server: %w", err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	// Authenticate with the SMTP server
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set the sender and recipient
	if err = c.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the email message
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to get email writer: %w", err)
	}

	message := []byte(fmt.Sprintf("Subject: Your OTP Code\n\nYour OTP code is: %s", otp))
	if _, err = w.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	// Quit the SMTP client
	if err = c.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP client: %w", err)
	}

	log.Println("Email sent successfully to:", to)
	return nil
}
