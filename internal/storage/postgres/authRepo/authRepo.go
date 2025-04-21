package authRepo

import (
	"context"
	"errors"
	"time"

	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage/postgres"
)

type authRepo struct {
	pool postgres.DBPool
}

func NewAuthRepo(pool postgres.DBPool) *authRepo {
	return &authRepo{
		pool: pool,
	}
}

func (r *authRepo) GetRoleIdByName(ctx context.Context, role string) (int, error) {
	query := `select id from role
				where name = $1`

	var id int
	err := r.pool.QueryRow(ctx, query, role).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *authRepo) CreateSession(ctx context.Context, user *models.User, sessionId string) (time.Time, error) {
	query := `insert into sessions(sessionId, userId, userRole, expireAt)
				values(
					$1, $2, $3, NOW() + INTERVAL '30 days'
				)
				returning expireAt`
	var expireAt time.Time
	err := r.pool.QueryRow(ctx, query, sessionId, user.Id, user.Role).Scan(&expireAt)

	return expireAt, err
}

func (r *authRepo) AddNewUser(ctx context.Context, user *models.User) error {
	query := `insert into users(email, password, roleId)
						values ($1, $2, $3)`

	roleId, err := r.GetRoleIdByName(ctx, user.Role)
	if err != nil {
		return errors.New("flag1#")
	}

	_, err = r.pool.Exec(ctx, query, user.Email, user.PasswordHash, roleId)
	if err != nil {
		return errors.New("flag2#")
	}
	return nil
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `select u.id, u.email, u.password, r.name
				from users u
				left join role r on u.roleid = r.id
				where u.email = $1`

	user := &models.User{}

	err := r.pool.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepo) UpdateSessionExpireTime(ctx context.Context, sessionId string) (time.Time, error) {
	query := `update sessions
				set expireAt = NOW() + INTERVAL '30 days'
				where sessionId = $1
				returning expireAt`

	var expireAt time.Time
	err := r.pool.QueryRow(ctx, query, sessionId).Scan(&expireAt)

	return expireAt, err
}

func (r *authRepo) GetSession(ctx context.Context, sessionId string) (*models.Session, error) {
	query := `select userId, userRole, expireAt
				from sessions
				where sessionId = $1`

	session := &models.Session{}

	err := r.pool.QueryRow(ctx, query, sessionId).Scan(&session.UserId, &session.UserRole, &session.ExpireAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}
