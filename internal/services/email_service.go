package services

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendPaymentSuccessEmailWithInvoice(
	to string,
	orderID string,
	amount int64,
	invoicePath string,
) error {

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Payment Successful - Wildan Store")

	body := `
Hi,

Your payment has been successfully processed.

Order ID : ` + orderID + `
Total    : Rp ` + strconv.FormatInt(amount, 10) + `
Status   : PAID

Please find your invoice attached.

Thank you,
Wildan Store
`
	m.SetBody("text/plain", body)

	// 📎 ATTACH INVOICE PDF
	if invoicePath != "" {
		m.Attach(invoicePath)
	}

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	return d.DialAndSend(m)
}
