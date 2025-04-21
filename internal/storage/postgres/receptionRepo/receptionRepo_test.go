package receptionRepo

import (
	"context"
	"errors"
	"fmt"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage/postgres"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockDbPool struct {
	mock.Mock
	postgres.DBPool
}

type mockDbTx struct {
	mock.Mock
	postgres.Tx
}

type mockRow struct {
	mock.Mock
}

func (m *mockDbPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	callArgs := m.Called(append([]interface{}{ctx, sql}, args...)...)
	return callArgs.Get(0).(pgx.Row)
}

func (m *mockDbTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	callArgs := m.Called(append([]interface{}{ctx, sql}, args...)...)
	return callArgs.Get(0).(pgx.Row)
}
func (m *mockDbTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	callArgs := m.Called(append([]interface{}{ctx, sql}, args...)...)
	return callArgs.Get(0).(pgconn.CommandTag), callArgs.Error(1)
}
func (m *mockDbPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	allArgs := m.Called(append([]interface{}{ctx, sql}, args...)...)
	callArgs := m.Called(allArgs...)
	return callArgs.Get(0).(pgconn.CommandTag), callArgs.Error(1)
}
func (m *mockDbPool) Begin(ctx context.Context) (postgres.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(postgres.Tx), args.Error(1)
}

func (m *mockRow) Scan(dest ...any) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *mockDbTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockDbTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestCloseReception(t *testing.T) {
	tests := []struct {
		name      string
		pvzId     uuid.UUID
		mockError error
	}{
		{
			name:      "valid request",
			pvzId:     uuid.New(),
			mockError: nil,
		},
		{
			name:      "invalid request",
			pvzId:     uuid.New(),
			mockError: errors.New("Bad request"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockPool := new(mockDbPool)
			repo := NewReceptionRepo(mockPool)

			mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, tt.mockError)

			err := repo.CloseReception(ctx, tt.pvzId)
			require.Equal(t, tt.mockError, err)

			mockPool.AssertExpectations(t)
		})
	}
}

func TestCreateReception(t *testing.T) {
	tests := []struct {
		name       string
		pvzId      uuid.UUID
		mockReturn *models.Reception
		mockError  error
	}{
		{
			name:       "invalid request",
			pvzId:      uuid.New(),
			mockReturn: nil,
			mockError:  errors.New("Bad request"),
		},
		{
			name:       "valid request",
			pvzId:      uuid.New(),
			mockReturn: &models.Reception{},
			mockError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			pgxRow := new(mockRow)
			mockPool := new(mockDbPool)
			repo := NewReceptionRepo(mockPool)

			mockPool.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(pgxRow)

			pgxRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.mockError)

			res, err := repo.CreateReception(ctx, tt.pvzId)

			require.Equal(t, err, tt.mockError)
			require.IsType(t, &models.Reception{}, res)

			mockPool.AssertExpectations(t)
		})
	}
}
func TestGetStatusNameById(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		mockReturn string
		mockError  error
	}{
		{
			name:       "invalid request",
			id:         -1,
			mockReturn: "",
			mockError:  errors.New("Bad request"),
		},
		{
			name:       "valid request",
			id:         1,
			mockReturn: "in_progress",
			mockError:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockPool := new(mockDbPool)
			repo := NewReceptionRepo(mockPool)
			pgxRow := new(mockRow)

			mockPool.On("QueryRow", ctx, mock.Anything, tt.id).Return(pgxRow)
			pgxRow.On("Scan", mock.Anything).Run(
				func(args mock.Arguments) {
					reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn))
				}).Return(tt.mockError)

			res, err := repo.GetStatusNameById(ctx, tt.id)
			require.Equal(t, err, tt.mockError)

			require.Equal(t, tt.mockReturn, res)

			mockPool.AssertExpectations(t)
		})
	}
}

func TestDeleteLastProductInReception(t *testing.T) {
	tests := []struct {
		name       string
		args       []uuid.UUID
		mockReturn []uuid.UUID
		mockError  error
	}{
		{
			name:       "invalid test",
			args:       []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			mockReturn: []uuid.UUID{uuid.New(), uuid.New()},
			mockError:  errors.New("error"),
		},
		{
			name:       "valid test",
			args:       []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			mockReturn: []uuid.UUID{uuid.New(), uuid.New()},
			mockError:  nil,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.mockReturn)
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockPool := new(mockDbPool)
			repo := NewReceptionRepo(mockPool)
			mockTx := new(mockDbTx)
			pgxRow1 := new(mockRow)
			pgxRow2 := new(mockRow)

			mockPool.On("Begin", ctx).Return(mockTx, tt.mockError)
			require.NotNil(t, mockTx)

			mockTx.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(pgxRow1)
			pgxRow1.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
				reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn[0]))
			}).Return(tt.mockError)

			mockTx.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(pgxRow2)
			pgxRow2.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
				reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn[1]))
			}).Return(tt.mockError)

			mockTx.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("DELETE 1"), tt.mockError)
			mockTx.On("Commit", ctx).Return(nil)

			mockTx.On("Rollback", ctx).Return(nil)

			err := repo.DeleteLastProductInReception(ctx, tt.args[0])

			require.Equal(t, err, tt.mockError)
			mockPool.AssertExpectations(t)
		})
	}
}

func TestAddProductToReception(t *testing.T) {
	tests := []struct {
		name       string
		argProd    *models.Product
		mockReturn *models.Product
		mockError  error
	}{
		{
			name:       "invalid test",
			argProd:    &models.Product{},
			mockReturn: &models.Product{ReceptionId: uuid.New(), Id: uuid.New(), AddedAt: time.Now()},
			mockError:  errors.New("error"),
		},
		{
			name:       "valid test",
			argProd:    &models.Product{},
			mockReturn: &models.Product{},
			mockError:  nil,
		},
	}

	for _, tt := range tests {
		fmt.Println(tt.mockReturn)
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockPool := new(mockDbPool)
			repo := NewReceptionRepo(mockPool)
			mockTx := new(mockDbTx)
			pgxRow1 := new(mockRow)
			pgxRow2 := new(mockRow)

			mockPool.On("Begin", ctx).Return(mockTx, tt.mockError)
			require.NotNil(t, mockTx)

			queryOpenReception := `select id
							from receptions
							where pvz_id = $1 and status_id = 1`
			mockTx.On("QueryRow", ctx, queryOpenReception, mock.Anything).Return(pgxRow1)
			pgxRow1.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
				reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn.ReceptionId))
			}).Return(tt.mockError)

			queryAddProduct := `insert into products(type_id)
						values ($1)
						returning id, added_at`
			mockTx.On("QueryRow", ctx, queryAddProduct, mock.Anything).Return(pgxRow2)
			pgxRow2.On("Scan", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn.Id))
				reflect.ValueOf(args[1]).Elem().Set(reflect.ValueOf(tt.mockReturn.AddedAt))
			}).Return(tt.mockError)

			mockTx.On("Exec", ctx, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("DELETE 1"), tt.mockError)

			mockTx.On("Commit", ctx).Return(nil)
			mockTx.On("Rollback", ctx).Return(nil)

			out, err := repo.AddProductToReception(ctx, tt.argProd, uuid.New())

			require.Equal(t, err, tt.mockError)
			if err == nil {
				mockPool.AssertExpectations(t)
				require.Equal(t, out.AddedAt, tt.mockReturn.AddedAt)
				require.Equal(t, out.Id, tt.mockReturn.Id)
				require.Equal(t, out.ReceptionId, tt.mockReturn.ReceptionId)
			}

		})
	}
}

func TestGetProductTypeIdByName(t *testing.T) {
	tests := []struct {
		name       string
		arg        string
		mockReturn int
		mockError  error
	}{
		{
			name:       "valid test",
			arg:        "Одежда",
			mockReturn: 1,
			mockError:  nil,
		},
		{
			name:       "invalid test",
			arg:        "еда",
			mockReturn: -1,
			mockError:  errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockPool := new(mockDbPool)
			repo := NewReceptionRepo(mockPool)
			pgxRow := new(mockRow)

			mockPool.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(pgxRow)
			pgxRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
				reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(tt.mockReturn))
			}).Return(tt.mockError)

			id, err := repo.GetProductTypeIdByName(ctx, tt.arg)
			require.Equal(t, err, tt.mockError)
			require.Equal(t, id, tt.mockReturn)
			mockPool.AssertExpectations(t)
		})
	}
}
