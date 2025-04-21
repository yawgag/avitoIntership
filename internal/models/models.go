package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	Id            uuid.UUID
	DateTime      time.Time
	PickupPointId uuid.UUID
	StatusId      int
}

type ReceptionAPI struct {
	Id            uuid.UUID `json:"id"`
	DateTime      time.Time `json:"dateTime"`
	PickupPointId uuid.UUID `json:"pvzId"`
	Status        string    `json:"status"`
}

type ProductAPI struct {
	Id          uuid.UUID  `json:"id"`
	AddedAt     time.Time  `json:"dateTime"`
	Type        string     `json:"type"`
	PvzId       *uuid.UUID `json:"pvzId,omitempty"`
	ReceptionId uuid.UUID  `json:"receptionId"`
}

type Product struct {
	Id          uuid.UUID
	AddedAt     time.Time
	TypeId      int
	ReceptionId uuid.UUID
}

type PvzFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PageLimit int
}

type PvzFilteredInfo struct {
	PvzID         uuid.UUID
	CityName      string
	RegDate       time.Time
	ReceptionID   uuid.UUID
	ReceptionTime time.Time
	ProductID     uuid.UUID
	AddedAt       time.Time
	ProductType   string
}

type ProductInfo struct {
	ID      uuid.UUID `json:"id"`
	AddedAt time.Time `json:"addedAt"`
	Type    string    `json:"type"`
}

type ReceptionInfo struct {
	ID       uuid.UUID     `json:"id"`
	DateTime time.Time     `json:"dateTime"`
	Products []ProductInfo `json:"products"`
}

type PvzInfo struct {
	ID         uuid.UUID       `json:"id"`
	CityName   string          `json:"city"`
	RegDate    time.Time       `json:"regDate"`
	Receptions []ReceptionInfo `json:"receptions"`
}

type PickupPoint struct {
	Id      uuid.UUID
	RegDate time.Time
	CityId  int
}

type PickupPointAPI struct {
	Id      uuid.UUID `json:"id"`
	RegDate time.Time `json:"registrationDate"`
	City    string    `json:"city"`
}

type User struct {
	Id           int     `json:"id"`
	Email        string  `json:"email"`
	Password     *string `json:"password"`
	PasswordHash *string `json:"passwordHash,omitempty"`
	Role         string  `json:"role"`
}

type Session struct {
	SessionId string    `json:"sessionId"`
	UserId    int       `json:"userId"`
	UserRole  string    `json:"userRole"`
	ExpireAt  time.Time `json:"expireAt"`
}

type AuthTokens struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	NewAccessToken  bool   `json:"-"`
	NewRefreshToken bool   `json:"-"`
}
