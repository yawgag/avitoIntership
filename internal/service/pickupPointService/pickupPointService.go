package pickupPointService

import (
	"context"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage"
)

type PickupPointService struct {
	PickupPointRepo storage.PickupPoint
}

func NewPickupPointService(pickupPointRepo storage.PickupPoint) *PickupPointService {
	return &PickupPointService{
		PickupPointRepo: pickupPointRepo,
	}
}

func (s *PickupPointService) Create(ctx context.Context, pickupPoint *models.PickupPoint) (*models.PickupPoint, error) {
	outPickupPoint, err := s.PickupPointRepo.Create(ctx, pickupPoint)
	if err != nil {
		return nil, err
	}
	return outPickupPoint, nil
}
