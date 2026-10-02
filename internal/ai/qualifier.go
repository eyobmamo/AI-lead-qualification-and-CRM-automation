package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	lead "lead-automation/internal/lead/model"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponseChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type Qualifier struct {
	url   string
	model string
}

func NewQualifier(url, model string) *Qualifier {
	return &Qualifier{
		url:   url,
		model: model,
	}
}

func (q *Qualifier) Qualify(
	ctx context.Context,
	l lead.Lead,
) (*lead.Qualification, error) {

	prompt := fmt.Sprintf(`
You are a lead qualification assistant.

Analyze this business lead:

Name: %s
Email: %s
Company: %s
Message: %s
Budget: %.2f
Timeline: %s

Return ONLY valid JSON in this exact format:

{
  "score": 0,
  "category": "HOT",
  "summary": "short explanation",
  "recommended_action": "what sales should do"
}

Rules:
- score must be between 0 and 100
- category must be HOT, WARM, or COLD
- Do not include markdown or extra text.
`,
		l.Name,
		l.Email,
		l.Company,
		l.Message,
		l.Budget,
		l.Timeline,
	)

	reqBody := GenerateRequest{
		Model:  q.model,
		Prompt: prompt,
		Stream: true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		q.url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var result strings.Builder

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		var chunk OllamaResponseChunk

		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			return nil, fmt.Errorf("invalid ollama response: %w", err)
		}

		result.WriteString(chunk.Response)

		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var qualification lead.Qualification

	if err := json.Unmarshal(
		[]byte(strings.TrimSpace(result.String())),
		&qualification,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to parse qualification: %w; response=%q",
			err,
			result.String(),
		)
	}

	if qualification.Score < 0 || qualification.Score > 100 {
		return nil, fmt.Errorf("invalid qualification score")
	}

	switch qualification.Category {
	case "HOT", "WARM", "COLD":
	default:
		return nil, fmt.Errorf(
			"invalid qualification category: %s",
			qualification.Category,
		)
	}

	return &qualification, nil
}