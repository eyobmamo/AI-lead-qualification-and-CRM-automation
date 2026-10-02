package lead

import (
	"context"

	"lead-automation/internal/ai"
	"lead-automation/internal/automation"
	lead "lead-automation/internal/lead/model"

)

type Service struct {
	repo *Repository
	ai   *ai.Qualifier
	n8n  *automation.N8N
}

func NewService(
	repo *Repository,
	ai *ai.Qualifier,
	n8n *automation.N8N,
) *Service {
	return &Service{
		repo: repo,
		ai:   ai,
		n8n:  n8n,
	}
}

func (s *Service) CreateLead(
	ctx context.Context,
	l *lead.Lead,
) (*lead.Lead, error) {

	// 1. Save lead.
	if err := s.repo.Create(ctx, l); err != nil {
		return nil, err
	}

	// 2. AI qualification will come here.
	// qualification, err := s.ai.Qualify(ctx, l)
	
	// 3. Send event to n8n.
	// err = s.n8n.SendLeadQualified(...)

	return l, nil
}