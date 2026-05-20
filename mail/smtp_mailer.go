package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailSender interface {
    SendVerificationEmail(toEmail, nickname, verificationURL string) error
}

type SMTPEmailSender struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
}

func NewSMTPEmailSender(host string, port int, username, password, from string) *SMTPEmailSender {
    return &SMTPEmailSender{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        From:     from,
    }
}

func (s *SMTPEmailSender) SendVerificationEmail(toEmail, nickname, verificationURL string) error {
    addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
    auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

    msg := strings.Join([]string{
        fmt.Sprintf("From: %s", s.From),
        fmt.Sprintf("To: %s", toEmail),
        "Subject: 请验证您的邮箱",
        "MIME-Version: 1.0",
        "Content-Type: text/plain; charset=UTF-8",
        "",
        fmt.Sprintf("你好 %s", nickname),
        "",
        "请点击下面链接验证你的邮箱：",
        verificationURL,
        "",
        "如果你没有发起该请求，请忽略此邮件。",
    }, "\r\n")

    // Gmail uses STARTTLS on port 587.
    client, err := smtp.Dial(addr)
    if err != nil {
        return err
    }
    defer client.Close()

    tlsConfig := &tls.Config{
        ServerName: s.Host,
    }
    if err := client.StartTLS(tlsConfig); err != nil {
        return err
    }

    if err := client.Auth(auth); err != nil {
        return err
    }

    if err := client.Mail(s.From); err != nil {
        return err
    }
    if err := client.Rcpt(toEmail); err != nil {
        return err
    }

    wc, err := client.Data()
    if err != nil {
        return err
    }
    defer wc.Close()

    if _, err := wc.Write([]byte(msg)); err != nil {
        return err
    }

    return nil
}
