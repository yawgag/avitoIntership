package authHandler

import (
	"encoding/json"
	"net/http"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/utils/errorsHandl"
	"time"
)

type authHandler struct {
	authService service.Auth
}

func NewAuthHandler(authService service.Auth) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var reqData models.User
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	mockUser := &models.User{
		Id:   -1,
		Role: reqData.Role,
	}

	tokens, err := h.authService.DummyLogin(r.Context(), mockUser)
	if err != nil {
		errorsHandl.SendJsonError(w, "bad request", http.StatusBadRequest)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   false, // TODO(changed for tests)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
	})
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorsHandl.SendJsonError(w, "flag1", http.StatusBadRequest)
		return
	}

	var reqData *models.User
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		errorsHandl.SendJsonError(w, "err1"+err.Error(), http.StatusBadRequest)
		return
	}
	// some validation
	if reqData.Email == "" || reqData.Password == "" || reqData.Role == "" {
		errorsHandl.SendJsonError(w, "Bad request. Bad values", http.StatusBadRequest)
		return
	}

	err := h.authService.Register(r.Context(), reqData)
	if err != nil {
		errorsHandl.SendJsonError(w, "err2"+err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusUnauthorized)
		return
	}

	var reqData *models.User
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusUnauthorized)
		return
	}

	tokens, err := h.authService.Login(r.Context(), reqData)
	if err != nil {
		errorsHandl.SendJsonError(w, "Wrong data", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		errorsHandl.SendJsonError(w, "bad request", http.StatusBadRequest)
	}
}

func (h *authHandler) IsSignedInMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("accessToken")
		if err != nil {
			errorsHandl.SendJsonError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		refreshTokenCookie, err := r.Cookie("refreshToken")
		if err != nil {
			errorsHandl.SendJsonError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokens := &models.AuthTokens{
			AccessToken:  accessTokenCookie.Value,
			RefreshToken: refreshTokenCookie.Value,
		}

		tokens, err = h.authService.HandleTokens(r.Context(), tokens)
		if err != nil {
			errorsHandl.SendJsonError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if tokens.NewRefreshToken {
			http.SetCookie(w, &http.Cookie{
				Name:     "refreshToken",
				Value:    tokens.RefreshToken,
				HttpOnly: true,
				Expires:  time.Now().Add(30 * 24 * time.Hour),
				Path:     "/",
			})
		}

		if tokens.NewAccessToken {
			http.SetCookie(w, &http.Cookie{
				Name:     "accessToken",
				Value:    tokens.AccessToken,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(15 * time.Minute),
				Path:     "/",
			})
		}
		next(w, r)
	})
}

func (h *authHandler) IsAvaliableRoleMiddleware(next http.HandlerFunc, avaliableRoles []string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("accessToken")
		if err != nil {
			errorsHandl.SendJsonError(w, "Unauthorized", http.StatusForbidden)
			return
		}
		tokens := &models.AuthTokens{
			AccessToken: accessTokenCookie.Value,
		}

		avaliable, err := h.authService.AvaliableForUser(tokens, avaliableRoles)
		if err != nil {
			errorsHandl.SendJsonError(w, "Unauthorized", http.StatusForbidden)
			return
		}

		if avaliable {
			next(w, r)
		} else {
			errorsHandl.SendJsonError(w, "Unauthorized", http.StatusForbidden)
		}
	})
}
