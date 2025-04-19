package receptionService

import (
	"context"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage"

	"github.com/google/uuid"
)

type ReceptionService struct {
	ReceptionRepo storage.Reception
}

func NewReceptionService(receptionRepo storage.Reception) *ReceptionService {
	return &ReceptionService{
		ReceptionRepo: receptionRepo,
	}
}

func (s *ReceptionService) CreateReception(ctx context.Context, pvzId uuid.UUID) (*models.ReceptionAPI, error) {
	reception, err := s.ReceptionRepo.CreateReception(ctx, pvzId)
	if err != nil {
		return nil, err
	}

	statusName, err := s.ReceptionRepo.GetStatusNameById(ctx, reception.StatusId)
	if err != nil {
		return nil, err
	}
	outReception := &models.ReceptionAPI{
		Id:            reception.Id,
		DateTime:      reception.DateTime,
		PickupPointId: reception.PickupPointId,
		Status:        statusName,
	}
	return outReception, nil
}

func (s *ReceptionService) AddProduct(ctx context.Context, productAPI *models.ProductAPI) (*models.ProductAPI, error) {
	typeId, err := s.ReceptionRepo.GetProductTypeIdByName(ctx, productAPI.Type)
	if err != nil {
		return nil, err
	}
	product := &models.Product{
		Id:      productAPI.Id,
		AddedAt: productAPI.AddedAt,
		TypeId:  typeId,
	}

	product, err = s.ReceptionRepo.AddProductToReception(ctx, product, *productAPI.PvzId)
	if err != nil {
		return nil, err
	}

	productAPI = &models.ProductAPI{
		Id:          product.Id,
		AddedAt:     product.AddedAt,
		Type:        productAPI.Type,
		ReceptionId: product.ReceptionId,
	}
	return productAPI, nil
}

// TODO: add validation
func (s *ReceptionService) DeleteLastProductInReception(ctx context.Context, pvzId uuid.UUID) error {
	err := s.ReceptionRepo.DeleteLastProductInReception(ctx, pvzId)
	if err != nil {
		return err
	}
	return nil
}

// TODO: i said "sad" about this function in repo, but this... much worse
func (s *ReceptionService) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	err := s.ReceptionRepo.CloseReception(ctx, pvzId)
	return err
}
