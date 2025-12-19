package services

type EmailJob struct {
	To          string `json:"to"`
	OrderID     string `json:"order_id"`
	Amount      int64  `json:"amount"`
	InvoicePath string `json:"invoice_path"`
}
