package receptionService

import (
	"context"
	"errors"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockReceptionRepo struct {
	mock.Mock
	storage.Reception
}

func (m *MockReceptionRepo) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	args := m.Called(ctx, pvzId)
	return args.Error(0)
}

func (m *MockReceptionRepo) CreateReception(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *MockReceptionRepo) GetStatusNameById(ctx context.Context, id int) (string, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockReceptionRepo) DeleteLastProductInReception(ctx context.Context, pvzId uuid.UUID) error {
	args := m.Called(ctx, pvzId)
	return args.Error(0)
}

func (m *MockReceptionRepo) GetProductTypeIdByName(ctx context.Context, name string) (int, error) {
	args := m.Called(ctx, name)
	return args.Int(0), args.Error(1)
}

func (m *MockReceptionRepo) AddProductToReception(ctx context.Context, product *models.Product, pvzId uuid.UUID) (*models.Product, error) {
	args := m.Called(ctx, product, pvzId)
	return args.Get(0).(*models.Product), args.Error(1)
}

func TestAddProduct_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockReceptionRepo)
	service := NewReceptionService(mockRepo)

	productType := "Электроника"
	typeId := 3
	pvzId := uuid.New()
	expectedProduct := &models.Product{
		Id:          uuid.New(),
		AddedAt:     time.Now(),
		TypeId:      typeId,
		ReceptionId: uuid.New(),
	}

	productAPI := &models.ProductAPI{
		Type:  productType,
		PvzId: &pvzId,
	}

	mockRepo.On("GetProductTypeIdByName", ctx, productType).Return(typeId, nil)
	mockRepo.On("AddProductToReception", ctx, mock.AnythingOfType("*models.Product"), pvzId).Return(expectedProduct, nil)

	result, err := service.AddProduct(ctx, productAPI)

	require.NoError(t, err)
	require.Equal(t, expectedProduct.Id, result.Id)
	require.Equal(t, expectedProduct.AddedAt, result.AddedAt)
	require.Equal(t, productType, result.Type)
	require.Equal(t, expectedProduct.ReceptionId, result.ReceptionId)

	mockRepo.AssertExpectations(t)
}

func TestCloseReception(t *testing.T) {
	tests := []struct {
		name      string
		arg       uuid.UUID
		mockError error
	}{
		{
			name:      "valid test",
			arg:       uuid.New(),
			mockError: nil,
		},
		{
			name:      "invalid test",
			arg:       uuid.New(),
			mockError: errors.New("Bad request"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := new(MockReceptionRepo)
			service := NewReceptionService(mockRepo)

			mockRepo.On("CloseReception", ctx, tt.arg).Return(tt.mockError)

			err := service.CloseReception(ctx, tt.arg)

			require.Equal(t, tt.mockError, err)
		})
	}
}

func TestDeleteLastProductInReception(t *testing.T) {
	tests := []struct {
		name      string
		arg       uuid.UUID
		mockError error
	}{
		{
			name:      "valid test",
			arg:       uuid.New(),
			mockError: nil,
		},
		{
			name:      "invalid test",
			arg:       uuid.New(),
			mockError: errors.New("Bad request"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := new(MockReceptionRepo)
			service := NewReceptionService(mockRepo)

			mockRepo.On("DeleteLastProductInReception", ctx, tt.arg).Return(tt.mockError)

			err := service.DeleteLastProductInReception(ctx, tt.arg)

			require.Equal(t, tt.mockError, err)
		})
	}
}

func TestCreateReception(t *testing.T) {
	tests := []struct {
		name       string
		arg        uuid.UUID
		mockReturn *models.ReceptionAPI
		mockError  error
	}{
		{
			name:       "valid test",
			arg:        uuid.New(),
			mockReturn: &models.ReceptionAPI{PickupPointId: uuid.New(), Status: "in_progress"},
			mockError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepo := new(MockReceptionRepo)
			service := NewReceptionService(mockRepo)

			mockRepo.On("CreateReception", ctx, mock.Anything).Return(&models.Reception{}, tt.mockError)
			mockRepo.On("GetStatusNameById", ctx, mock.Anything).Return(tt.mockReturn.Status, tt.mockError)

			out, err := service.CreateReception(ctx, tt.mockReturn.PickupPointId)
			require.Equal(t, tt.mockError, err)
			if tt.mockError == nil {
				require.Equal(t, tt.mockReturn, out)
				mockRepo.AssertExpectations(t)
			}

		})
	}
}
