package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type InvoiceData struct {
	OrderID     string
	Email       string
	GameName    string
	PackageName string
	Amount      int64
	PaidAt      time.Time
}

func saveEmbeddedLogo() (string, error) {
	path := "tmp_logo.png"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return path, os.WriteFile(path, Logo, 0644)
}

func GenerateInvoicePDF(data InvoiceData) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// ===============================
	// HEADER (LOGO + BRAND)
	// ===============================
	logoPath, _ := saveEmbeddedLogo()
	pdf.Image(logoPath, 10, 10, 40, 0, false, "", 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(120, 15)
	pdf.Cell(40, 10, "INVOICE")

	pdf.Ln(30)

	// ===============================
	// STORE INFO
	// ===============================
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 6, "Wildan Store")
	pdf.Ln(5)
	pdf.Cell(40, 6, "Email: wildanhanifabdillah27@gmail.com")
	pdf.Ln(8)

	// ===============================
	// INVOICE META
	// ===============================
	pdf.Cell(40, 6, fmt.Sprintf("Invoice ID : %s", data.OrderID))
	pdf.Ln(5)
	pdf.Cell(40, 6, fmt.Sprintf("Paid Date  : %s", data.PaidAt.Format("02 Jan 2006 15:04")))
	pdf.Ln(8)

	// ===============================
	// CUSTOMER
	// ===============================
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(40, 6, "Bill To:")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 6, data.Email)
	pdf.Ln(10)

	// ===============================
	// TABLE HEADER
	// ===============================
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(240, 240, 240)

	pdf.CellFormat(90, 8, "Description", "1", 0, "", true, 0, "")
	pdf.CellFormat(30, 8, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Price", "1", 1, "R", true, 0, "")

	// ===============================
	// TABLE CONTENT
	// ===============================
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(90, 8, data.GameName+" - "+data.PackageName, "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 8, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("Rp %d", data.Amount), "1", 1, "R", false, 0, "")

	// ===============================
	// TOTAL
	// ===============================
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(120, 8, "Total", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("Rp %d", data.Amount), "1", 1, "R", false, 0, "")

	// ===============================
	// FOOTER
	// ===============================
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, "Thank you for your purchase at Wildan Store")
	pdf.Ln(5)
	pdf.Cell(0, 6, "This invoice is generated automatically and valid without signature")

	// ===============================
	// SAVE FILE
	// ===============================
	dir := "invoices"
	_ = os.MkdirAll(dir, os.ModePerm)

	filePath := filepath.Join(dir, data.OrderID+".pdf")
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
