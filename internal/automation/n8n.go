package automation

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type N8N struct {
	webhookURL string
}

func NewN8N(webhookURL string) *N8N {
	return &N8N{
		webhookURL: webhookURL,
	}
}

type LeadQualifiedEvent struct {
	Event    string `json:"event"`
	LeadID   int64  `json:"lead_id"`
	Score    int    `json:"score"`
	Category string `json:"category"`
}

func (n *N8N) SendLeadQualified(
	ctx context.Context,
	event LeadQualifiedEvent,
) error {

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.webhookURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil // we'll improve error handling later
	}

	return nil
}