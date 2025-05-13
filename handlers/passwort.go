package handlers

import (
	"fmt"
	"net/http"
	"os"
	"net/smtp"
	"regexp"
	"zeiterfassung-backend/config"
	"zeiterfassung-backend/models"

	"github.com/gin-gonic/gin"
)

// POST /api/passwort-vergessen
func PasswortVergessenHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil || request.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Eingabe"})
		return
	}

	// Email-Format validieren
	if !isValidEmail(request.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültiges E-Mail-Format"})
		return
	}

	// Prüfen, ob der Nutzer existiert
	var nutzer models.Nutzer
	result := config.DB.Where("email = ?", request.Email).First(&nutzer)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Es wurde kein Konto mit dieser E-Mail-Adresse gefunden"})
		return
	}

	// E-Mail an Admin senden
	if err := sendeAdminMail(request.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "E-Mail konnte nicht gesendet werden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Der Administrator wurde benachrichtigt"})
}

// E-Mail-Format validieren mit regulärem Ausdruck
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func sendeAdminMail(userEmail string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	to := os.Getenv("ADMIN_EMAIL")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || pass == "" || to == "" || smtpHost == "" || smtpPort == "" {
		fmt.Println("[WARNUNG] SMTP-Umgebungsvariablen fehlen – simuliere Mailversand")
		fmt.Printf("An: %s\nBetreff: Passwort-Zurücksetzen\nInhalt:\nEin Nutzer möchte ein neues Passwort: %s\n", to, userEmail)
		return nil // kein Fehler, nur kein echter Versand
	}

	subject := "Passwort-Zurücksetzen angefragt"
	body := fmt.Sprintf("Ein Nutzer möchte ein neues Passwort:\n\n%s", userEmail)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		[]byte(msg),
	)

	return err
}