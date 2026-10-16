package automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	lead "lead-automation/internal/lead/model"
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
	event lead.Qualification,
) error {

	body, err := json.Marshal(event)
	if err != nil {
			fmt.Println("__________sendleadQualified email error___________")

		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.webhookURL,
		bytes.NewReader(body),
	)
	fmt.Printf("webhookUrl: %v",n.webhookURL)
	if err != nil {
			fmt.Println("__________request dosent comen correclty to after call___________")

		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("__________callein http error___________")

		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {

		fmt.Println("-----------------------------------------")
		fmt.Println(resp) 
		fmt.Println("-----------------------------------------")

		fmt.Println("__________error status___________")

		return nil // we'll improve error handling later
	}

	fmt.Println("__________web hook is successfuly___________")
	return nil
}