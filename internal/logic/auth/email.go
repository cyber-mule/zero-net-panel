package auth

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/config"
)

func sendAuthEmail(logger logx.Logger, cfg config.AuthEmailConfig, to, subject, body string) error {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if provider == "" || provider == "log" {
		logger.Infof("mail: to=%s subject=%s body=%s", to, subject, body)
		return nil
	}
	if provider != "smtp" {
		return fmt.Errorf("auth: unsupported email provider %q", cfg.Provider)
	}

	from := strings.TrimSpace(cfg.From)
	if from == "" {
		return fmt.Errorf("auth: email from address is required")
	}

	host := strings.TrimSpace(cfg.SMTP.Host)
	if host == "" {
		return fmt.Errorf("auth: smtp host is required")
	}
	port := cfg.SMTP.Port
	if port == 0 {
		port = 587
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	message := buildMailMessage(from, to, subject, body)

	var auth smtp.Auth
	if cfg.SMTP.Username != "" || cfg.SMTP.Password != "" {
		auth = smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, host)
	}

	if cfg.SMTP.UseTLS {
		return sendMailTLS(addr, host, auth, from, []string{to}, []byte(message))
	}

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
}

// SendAuthEmail exposes email sending for auth-related workflows.
func SendAuthEmail(logger logx.Logger, cfg config.AuthEmailConfig, to, subject, body string) error {
	return sendAuthEmail(logger, cfg, to, subject, body)
}

func buildMailMessage(from, to, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}

func sendMailTLS(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write(msg); err != nil {
		_ = writer.Close()
		return err
	}
	return writer.Close()
}
