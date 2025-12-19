package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
)

type SnapRequest struct {
	TransactionDetails struct {
		OrderID     string `json:"order_id"`
		GrossAmount int64  `json:"gross_amount"`
	} `json:"transaction_details"`
	CustomerDetails struct {
		Email string `json:"email"`
	} `json:"customer_details"`
}

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func CreateSnap(orderID string, amount int64, email string) (*SnapResponse, error) {
	reqBody := SnapRequest{}
	reqBody.TransactionDetails.OrderID = orderID
	reqBody.TransactionDetails.GrossAmount = amount
	reqBody.CustomerDetails.Email = email

	payload, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("MIDTRANS_BASE_URL")+"/snap/v1/transactions",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, err
	}

	// Basic Auth: server_key:
	auth := base64.StdEncoding.EncodeToString(
		[]byte(os.Getenv("MIDTRANS_SERVER_KEY") + ":"),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var snapRes SnapResponse
	if err := json.NewDecoder(res.Body).Decode(&snapRes); err != nil {
		return nil, err
	}

	return &snapRes, nil
}
