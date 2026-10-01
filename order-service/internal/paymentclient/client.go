package paymentclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MOMON8798/Event-Driven.git/internal/service"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type createPaymentRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

type paymentResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (c *Client) CreatePayment(ctx context.Context, orderID string, amount float64) (*service.PaymentResult, error) {
	reqBody, err := json.Marshal(createPaymentRequest{OrderID: orderID, Amount: amount})
	if err != nil {
		return nil, fmt.Errorf("paymentclient: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/payments/", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("paymentclient: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paymentclient: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("paymentclient: unexpected status code %d", resp.StatusCode)
	}

	var payment paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		return nil, fmt.Errorf("paymentclient: decode response: %w", err)
	}

	return &service.PaymentResult{ID: payment.ID, Status: payment.Status}, nil
}
