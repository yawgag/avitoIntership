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

func (s *PickupPointService) Create(ctx context.Context, pickupPointAPI *models.PickupPointAPI) (*models.PickupPointAPI, error) {
	cityId, err := s.PickupPointRepo.GetCityIdByName(ctx, pickupPointAPI.City)

	pickupPoint := &models.PickupPoint{
		CityId: cityId,
	}

	pickupPoint, err = s.PickupPointRepo.Create(ctx, pickupPoint)
	if err != nil {
		return nil, err
	}

	outPickupPoint := &models.PickupPointAPI{
		Id:      pickupPoint.Id,
		RegDate: pickupPoint.RegDate,
		City:    pickupPointAPI.City,
	}

	return outPickupPoint, nil
}
