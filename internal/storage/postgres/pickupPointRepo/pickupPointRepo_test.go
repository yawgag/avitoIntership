package pickupPointRepo

import (
	"context"
	"errors"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage/postgres"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockDbPool struct {
	mock.Mock
	postgres.DBPool
}

type mockRow struct {
	mock.Mock
}

func (m *mockDbPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	callArgs := m.Called(append([]interface{}{ctx, sql}, args...)...)
	return callArgs.Get(0).(pgx.Row)
}

func (m *mockRow) Scan(dest ...any) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func TestGetCityIdByName(t *testing.T) {
	tests := []struct {
		name       string
		arg        string
		mockReturn int
		mockError  error
	}{
		{
			name:       "valid test",
			arg:        "Москва",
			mockReturn: 1,
			mockError:  nil,
		},
		{
			name:       "invalid test",
			arg:        "Омск",
			mockReturn: -1,
			mockError:  errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockPool := new(mockDbPool)
			repo := NewPickupPointRepo(mockPool)
			pgxRow := new(mockRow)

			mockPool.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(pgxRow)
			pgxRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
				reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn))
			}).Return(tt.mockError)

			id, err := repo.GetCityIdByName(ctx, tt.arg)
			require.Equal(t, err, tt.mockError)
			require.Equal(t, id, tt.mockReturn)
			mockPool.AssertExpectations(t)
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		arg        *models.PickupPoint
		mockReturn *models.PickupPoint
		mockError  error
	}{
		{
			name:       "valid test",
			arg:        &models.PickupPoint{CityId: 1},
			mockReturn: &models.PickupPoint{CityId: 0, RegDate: time.Now(), Id: uuid.New()},
			mockError:  nil,
		},
		{
			name:       "invalid test",
			arg:        &models.PickupPoint{},
			mockReturn: &models.PickupPoint{},
			mockError:  errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockPool := new(mockDbPool)
			repo := NewPickupPointRepo(mockPool)
			pgxRow := new(mockRow)

			mockPool.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(pgxRow)
			pgxRow.On("Scan", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn.Id))
				reflect.ValueOf(args[1]).Elem().Set(reflect.ValueOf(tt.mockReturn.RegDate))
			}).Return(tt.mockError)

			out, err := repo.Create(ctx, tt.arg)
			require.Equal(t, err, tt.mockError)
			if err == nil {
				require.Equal(t, out, tt.mockReturn)
			}

			mockPool.AssertExpectations(t)
		})
	}
}
