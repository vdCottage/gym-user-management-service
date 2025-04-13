package auth

import (
	"fmt"
	"net/smtp"

	"github.com/yourname/fitness-platform/config"
	"github.com/yourname/fitness-platform/pkg/logger"
)

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string, cfg *config.Config) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", cfg.SMTP.User, cfg.SMTP.Password, cfg.SMTP.Host)

	// Set up email headers
	headers := make(map[string]string)
	headers["From"] = cfg.SMTP.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build email message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Send email
	addr := fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port)
	err := smtp.SendMail(addr, auth, cfg.SMTP.From, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Create a logger instance
	log := logger.NewLogger("info", "json", "")
	log.Info("Email sent successfully", map[string]interface{}{
		"to":      to,
		"subject": subject,
	})

	return nil
}

// SendVerificationEmail sends a verification email with a token
func SendVerificationEmail(to, token string, cfg *config.Config) error {
	subject := "Verify Your Email"

	// Create verification link
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", cfg.App.URL, token)

	// Create email body
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Email Verification</h2>
				<p>Please click the link below to verify your email address:</p>
				<p><a href="%s">Verify Email</a></p>
				<p>If you did not request this verification, please ignore this email.</p>
				<p>This link will expire in %s.</p>
				<p>Best regards,<br>%s Team</p>
			</body>
		</html>
	`, verificationLink, cfg.JWT.Expiration, cfg.App.Name)

	return SendEmail(to, subject, body, cfg)
}

// SendPasswordResetEmail sends a password reset email with a token
func SendPasswordResetEmail(to, token string, cfg *config.Config) error {
	subject := "Reset Your Password"

	// Create reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", cfg.App.URL, token)

	// Create email body
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Password Reset</h2>
				<p>Please click the link below to reset your password:</p>
				<p><a href="%s">Reset Password</a></p>
				<p>If you did not request this password reset, please ignore this email.</p>
				<p>This link will expire in %s.</p>
				<p>Best regards,<br>%s Team</p>
			</body>
		</html>
	`, resetLink, cfg.JWT.Expiration, cfg.App.Name)

	return SendEmail(to, subject, body, cfg)
}
