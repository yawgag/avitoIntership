package pickupPointService

import (
	"context"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage"

	"github.com/google/uuid"
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

func (s *PickupPointService) GetInfo(ctx context.Context, filter *models.PvzFilter) ([]models.PvzInfo, error) {
	info, err := s.PickupPointRepo.GetFilteredInfo(ctx, filter)
	if err != nil {
		return nil, err
	}
	pvzMap := make(map[uuid.UUID]*models.PvzInfo)

	for _, item := range info {
		pvz, exists := pvzMap[item.PvzID]
		if !exists {
			pvz = &models.PvzInfo{
				ID:       item.PvzID,
				CityName: item.CityName,
				RegDate:  item.RegDate,
			}
			pvzMap[item.PvzID] = pvz
		}

		var rec *models.ReceptionInfo
		for i := range pvz.Receptions {
			if pvz.Receptions[i].ID == item.ReceptionID {
				rec = &pvz.Receptions[i]
				break
			}
		}
		if rec == nil {
			rec = &models.ReceptionInfo{
				ID:       item.ReceptionID,
				DateTime: item.ReceptionTime,
			}
			pvz.Receptions = append(pvz.Receptions, *rec)
			rec = &pvz.Receptions[len(pvz.Receptions)-1]
		}

		product := models.ProductInfo{
			ID:      item.ProductID,
			AddedAt: item.AddedAt,
			Type:    item.ProductType,
		}
		rec.Products = append(rec.Products, product)
	}

	result := make([]models.PvzInfo, 0, len(pvzMap))
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}
	return result, nil
}
