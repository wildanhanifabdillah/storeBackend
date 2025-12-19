package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPaymentSuccessEmail(to string, orderID string, amount int64) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	subject := "Payment Successful - Wildan Store"

	body := fmt.Sprintf(`
Hi,

Your payment has been successfully processed.

Order ID : %s
Total    : Rp %d
Status   : PAID

Thank you for using Wildan Store.
`, orderID, amount)

	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body,
	)

	auth := smtp.PlainAuth("", from, password, host)
	addr := host + ":" + port

	return smtp.SendMail(addr, auth, from, []string{to}, message)
}
