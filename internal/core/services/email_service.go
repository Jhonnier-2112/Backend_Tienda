package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
)

type emailService struct {
	host string
	port string
	user string
	pass string
	from string
	repo ports.EmailLogRepository
}

func NewEmailService(repo ports.EmailLogRepository) ports.EmailService {
	return &emailService{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
		from: os.Getenv("SMTP_FROM"),
		repo: repo,
	}
}

func (s *emailService) sendEmail(to, subject, body string) error {
	// If SMTP is not configured, log it and return nil to avoid crashing (for now)
	if s.host == "" || s.user == "" {
		log.Printf("[EMAIL MOCK/LOG] To: %s | Subject: %s | Body: %s", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n%s",
		s.from, to, subject, body))

	logEntry := &domain.EmailLog{
		ToAddress: to,
		Subject:   subject,
		Body:      body,
		Status:    "sent",
	}

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	if err != nil {
		log.Printf("[ERROR] Failed to send email to %s: %v", to, err)
		logEntry.Status = "failed"
		logEntry.Error = err.Error()
		s.repo.Create(logEntry)
		return err
	}

	s.repo.Create(logEntry)
	log.Printf("[SUCCESS] Email sent to %s", to)
	return nil
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	setupURL := fmt.Sprintf("http://localhost:3000/setup-password?token=%s", token)
	body := fmt.Sprintf(`
		<h1>Bienvenido a Nuestra Tienda</h1>
		<p>Por favor, haz clic en el siguiente enlace para verificar tu correo electrónico y configurar tu contraseña:</p>
		<a href="%s">Configurar Contraseña</a>
		<p>Si no creaste esta cuenta, puedes ignorar este mensaje.</p>
	`, setupURL)

	return s.sendEmail(to, "Verifica tu cuenta y configura tu contraseña", body)
}

func (s *emailService) SendPasswordResetEmail(to, token string) error {
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	body := fmt.Sprintf(`
		<h1>Recuperación de Contraseña</h1>
		<p>Hemos recibido una solicitud para restablecer tu contraseña. Haz clic en el siguiente enlace:</p>
		<a href="%s">Restablecer Contraseña</a>
		<p>Este enlace expirará en 1 hora.</p>
		<p>Si no solicitaste esto, ignora este mensaje.</p>
	`, resetURL)

	return s.sendEmail(to, "Recupera tu contraseña", body)
}
