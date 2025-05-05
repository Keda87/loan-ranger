package mailer

import "context"

type SendEmail struct {
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	ToEmail     string   `json:"to_email,omitempty"`
	ToListEmail []string `json:"to_list_email,omitempty"`
}

type EmailClientInterface interface {
	SendEmail(ctx context.Context, data SendEmail) error
}
