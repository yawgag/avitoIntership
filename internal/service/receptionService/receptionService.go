package receptionService

import (
	"context"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage"
)

type ReceptionService struct {
	ReceptionRepo storage.Reception
}

func NewReceptionService(receptionRepo storage.Reception) *ReceptionService {
	return &ReceptionService{
		ReceptionRepo: receptionRepo,
	}
}

func (s *ReceptionService) CreateReception(ctx context.Context, pvzId int) (*models.Reception, error) {
	// outReception, err := s.ReceptionRepo.CreateReception(ctx, pvzId)

	// statusName, err := s.ReceptionRepo.GetStatusNameById(ctx, outReception.Status)
	return nil, nil
}
