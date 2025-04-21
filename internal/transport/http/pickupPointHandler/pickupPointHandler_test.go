package pickupPointHandler

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

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPickupPoint struct {
	service.PickupPoint
	mock.Mock
}

func (m *mockPickupPoint) Create(ctx context.Context, pickupPoint *models.PickupPointAPI) (*models.PickupPointAPI, error) {
	args := m.Called(ctx, pickupPoint)
	return args.Get(0).(*models.PickupPointAPI), args.Error(1)
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		requestBody  interface{}
		mockReturn   *models.PickupPointAPI
		mockError    error
		answerStatus int
	}{
		{
			name:         "valid request",
			contentType:  "application/json",
			requestBody:  &models.PickupPointAPI{City: "Казань"},
			mockReturn:   &models.PickupPointAPI{City: "Казань"},
			mockError:    nil,
			answerStatus: http.StatusOK,
		},
		{
			name:         "wrong content type",
			contentType:  "text/plain",
			requestBody:  &models.PickupPointAPI{City: "Казань"},
			mockReturn:   nil,
			mockError:    nil,
			answerStatus: http.StatusBadRequest,
		},
		{
			name:         "invalid city",
			contentType:  "application/json",
			requestBody:  &models.PickupPointAPI{City: "Омск"},
			mockReturn:   &models.PickupPointAPI{},
			mockError:    errors.New("Bad request"),
			answerStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mockPickupPoint)
			handler := NewPickupPointHandler(mockService)

			requestBodyJson := []byte{}
			if tt.requestBody != nil {
				requestBodyJson, _ = json.Marshal(tt.requestBody)
			}

			httpRequset := httptest.NewRequest("POST", "/pvz", bytes.NewBuffer(requestBodyJson))
			httpRequset.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()

			if tt.mockReturn != nil || tt.mockError != nil {
				mockService.On("Create", mock.Anything, tt.requestBody.(*models.PickupPointAPI)).Return(tt.mockReturn, tt.mockError)
			}

			handler.Create(rec, httpRequset)
			require.Equal(t, tt.answerStatus, rec.Code)

			if rec.Code == http.StatusOK {
				var response models.PickupPointAPI
				err := json.NewDecoder(rec.Body).Decode(&response)
				require.NoError(t, err)

				require.Equal(t, tt.mockReturn.City, response.City)
			}

			if tt.mockReturn != nil || tt.mockError != nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}
