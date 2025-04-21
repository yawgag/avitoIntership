package receptionHandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockReceptionService struct {
	mock.Mock
	service.Reception
}

func (m *mockReceptionService) CreateReception(ctx context.Context, pvzId uuid.UUID) (*models.ReceptionAPI, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(*models.ReceptionAPI), args.Error(1)
}

func (m *mockReceptionService) AddProduct(ctx context.Context, productAPI *models.ProductAPI) (*models.ProductAPI, error) {
	args := m.Called(ctx, productAPI)
	return args.Get(0).(*models.ProductAPI), args.Error(1)
}

func (m *mockReceptionService) DeleteLastProductInReception(ctx context.Context, pvzId uuid.UUID) error {
	args := m.Called(ctx, pvzId)
	return args.Error(0)
}

func (m *mockReceptionService) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	args := m.Called(ctx, pvzId)
	return args.Error(0)
}

func TestCreateReception(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		requestBody  interface{}
		mockReturn   *models.ReceptionAPI
		answerStatus int
		mockError    error
	}{
		{
			name:         "Invalid content Type",
			contentType:  "text/plain",
			requestBody:  nil,
			answerStatus: http.StatusBadRequest,
		},
		{
			name:         "Empty json",
			contentType:  "application/json",
			requestBody:  "invalid json",
			answerStatus: http.StatusBadRequest,
		},
		{
			name:         "valid request",
			contentType:  "application/json",
			requestBody:  models.ReceptionAPI{PickupPointId: uuid.New()},
			mockReturn:   &models.ReceptionAPI{Status: "in_progress"},
			answerStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mockReceptionService)
			handler := NewReceptionHandler(mockService)

			requestBodyJson := []byte{}
			if tt.requestBody != nil {
				requestBodyJson, _ = json.Marshal(tt.requestBody)
			}

			httpRequset := httptest.NewRequest("POST", "/receptions", bytes.NewBuffer(requestBodyJson))
			httpRequset.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()

			if tt.mockReturn != nil {
				pvzId := tt.requestBody.(models.ReceptionAPI).PickupPointId
				mockService.On("CreateReception", mock.Anything, pvzId).Return(tt.mockReturn, tt.mockError)
			}
			handler.CreateReception(rec, httpRequset)

			require.Equal(t, tt.answerStatus, rec.Code)

			if rec.Code == http.StatusOK {
				var response models.ReceptionAPI
				err := json.NewDecoder(rec.Body).Decode(&response)
				require.NoError(t, err)
				require.Equal(t, "in_progress", response.Status)
			}

			if tt.mockReturn != nil || tt.mockError != nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestAddProduct(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		requestBody  interface{}
		mockReturn   *models.ProductAPI
		mockError    error
		answerStatus int
	}{
		{
			name:         "valid request",
			contentType:  "application/json",
			requestBody:  &models.ProductAPI{Type: "Одежда", ReceptionId: uuid.New()},
			mockReturn:   &models.ProductAPI{Type: "Одежда", ReceptionId: uuid.New()},
			mockError:    nil,
			answerStatus: http.StatusOK,
		},
		{
			name:         "Invalid content Type",
			contentType:  "text/plain",
			requestBody:  nil,
			answerStatus: http.StatusBadRequest,
		},
		{
			name:         "Empty json",
			contentType:  "application/json",
			requestBody:  "invalid json",
			mockError:    nil,
			answerStatus: http.StatusBadRequest,
		},
		{
			name:         "Wrong product type",
			contentType:  "application/json",
			requestBody:  &models.ProductAPI{Type: "еда", ReceptionId: uuid.New()},
			mockReturn:   &models.ProductAPI{},
			mockError:    errors.New("bad request"),
			answerStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mockReceptionService)
			handler := NewReceptionHandler(mockService)

			requestBodyJson := []byte{}
			if tt.requestBody != nil {
				requestBodyJson, _ = json.Marshal(tt.requestBody)
			}

			httpRequset := httptest.NewRequest("POST", "/products", bytes.NewBuffer(requestBodyJson))
			httpRequset.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()

			if tt.mockReturn != nil {
				mockService.On("AddProduct", mock.Anything, tt.requestBody.(*models.ProductAPI)).Return(tt.mockReturn, tt.mockError)
			}

			handler.AddProduct(rec, httpRequset)

			require.Equal(t, rec.Code, tt.answerStatus)

			if rec.Code == http.StatusOK {
				var response models.ProductAPI
				err := json.NewDecoder(rec.Body).Decode(&response)
				require.NoError(t, err)
				require.Equal(t, response.ReceptionId, tt.mockReturn.ReceptionId)
				require.Equal(t, response.Type, tt.mockReturn.Type)
			}

			if tt.mockReturn != nil || tt.mockError != nil {
				mockService.AssertExpectations(t)
			}

		})
	}
}

func TestDeleteLastProductInReception(t *testing.T) {
	tests := []struct {
		name         string
		pvzId        string
		mockError    error
		answerStatus int
	}{
		// only 2 cases exist
		{
			name:         "valid uuid",
			pvzId:        uuid.New().String(),
			mockError:    nil,
			answerStatus: http.StatusOK,
		},
		{
			name:         "invalid uuid",
			pvzId:        "invalid_uuid",
			mockError:    errors.New("Bad request"),
			answerStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mockReceptionService)
			handler := NewReceptionHandler(mockService)

			httpRequset := httptest.NewRequest("POST", "/pvz/"+tt.pvzId+"/delete_last_product", nil)
			httpRequset = mux.SetURLVars(httpRequset, map[string]string{"pvzId": tt.pvzId})
			rec := httptest.NewRecorder()
			if tt.mockError == nil {
				parsedId, _ := uuid.Parse(tt.pvzId)
				mockService.On("DeleteLastProductInReception", mock.Anything, parsedId).Return(tt.mockError)
			}
			handler.DeleteLastProduct(rec, httpRequset)

			require.Equal(t, rec.Code, tt.answerStatus)
			if tt.mockError == nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestCloseReception(t *testing.T) {
	tests := []struct {
		name         string
		pvzId        string
		mockError    error
		answerStatus int
	}{
		// only 2 cases exist
		{
			name:         "valid uuid",
			pvzId:        uuid.New().String(),
			mockError:    nil,
			answerStatus: http.StatusOK,
		},
		{
			name:         "invalid uuid",
			pvzId:        "invalid_uuid",
			mockError:    errors.New("Bad request"),
			answerStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mockReceptionService)
			handler := NewReceptionHandler(mockService)

			httpRequset := httptest.NewRequest("POST", "/pvz/"+tt.pvzId+"/close_last_reception", nil)
			httpRequset = mux.SetURLVars(httpRequset, map[string]string{"pvzId": tt.pvzId})
			rec := httptest.NewRecorder()
			if tt.mockError == nil {
				parsedId, _ := uuid.Parse(tt.pvzId)
				mockService.On("CloseReception", mock.Anything, parsedId).Return(tt.mockError)
			}
			handler.CloseReception(rec, httpRequset)

			require.Equal(t, rec.Code, tt.answerStatus)
			if tt.mockError == nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}
