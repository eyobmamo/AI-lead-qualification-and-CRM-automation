package lead

import "time"

type Lead struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	Message   string    `json:"message"`
	Budget    float64   `json:"budget"`
	Timeline  string    `json:"timeline"`
	CreatedAt time.Time `json:"created_at"`
}

type Qualification struct {
	Score              int     `json:"score"`
	Category           string  `json:"category"`
	Summary            string  `json:"summary"`
	RecommendedAction  string  `json:"recommended_action"`
}