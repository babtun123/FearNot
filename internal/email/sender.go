package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	From     string
	Password string
	To       string
	SMTPHost string
	SMTPPort string
}

type Sender struct {
	logger *log.Logger
}

func NewSender(logger *log.Logger) *Sender {
	return &Sender{
		logger: logger,
	}
}

// sendEmail prepares the email setup to be sent
func (s *Sender) sendEmail(config Config, subject, body string) error {

	// Build the email message in RFC 822 format
	// This format is required by SMTP servers

	message := []byte(
		"From: " + config.From + "\r\n" +
			"To: " + config.To + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
			"\r\n" +
			body + "\r\n")

	auth := smtp.PlainAuth("", config.From, config.Password, config.SMTPHost)

	err := smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort,
		auth,
		config.From,
		[]string{config.To},
		message,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// sendDailyVerse prepares the message to be sent
func (s *Sender) sendDailyVerse(config Config, verse string) error {
	subject := fmt.Sprintf("Daily Bible Verse - %s", time.Now().Format("Monday, January 2, 2006"))
	body := fmt.Sprintf(`Hello!

Here is your daily Bible verse, a reminder from God to not be afraid:

	"%s"

---
Sent with love from your Daily Verse app

To stop receiving these emails, please reply to this message.`, verse)

	return s.sendEmail(config, subject, body)
}

func GenerateEmail(logger *log.Logger, emailList []string, verseOfTheDay, ScriptureOfTheDay string) error {
	err := godotenv.Load()
	if err != nil {
		logger.Println("Error loading .env file")
		return fmt.Errorf("error loading .env file: %w", err)
	}
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailSender := os.Getenv("EMAIL_SENDER")

	em := NewSender(logger)

	verse := ScriptureOfTheDay + " - " + verseOfTheDay + " [KJV]"

	for _, email := range emailList {
		config := Config{
			From:     emailSender,
			Password: emailPassword,
			To:       email,
			SMTPHost: "smtp.gmail.com",
			SMTPPort: "587",
		}
		fmt.Printf("Sending daily verse email to: '%s'\n", email)
		err := em.sendDailyVerse(config, verse)

		if err != nil {
			logger.Println("Error sending daily verse email to")
			fmt.Printf("❌ Error sending email to '%s': %v\n", email, err)
		} else {
			logger.Println("Daily verse email sent successfully")
			fmt.Printf("✅ Email sent successfully to: '%s'\n!", email)
		}
	}

	return nil
}
