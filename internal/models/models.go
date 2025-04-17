package models

import (
	"time"
)

type Reception struct {
	Id            int       `json:"id"`
	DateTime      time.Time `json:"dateTime"`
	PickupPointId int       `json:"pvzId"`
	Status        int       `json:"status"`
}

type Product struct {
	Id      int       `json:"id"`
	AddedAt time.Time `json:"addedAt"`
	Type    string    `json:"type"`
}

type PickupPoint struct {
	Id      int       `json:"id"`
	RegDate time.Time `json:"registrationDate"`
	City    string    `json:"city"`
}

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordHash"`
	Role         string `json:"role"`
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
