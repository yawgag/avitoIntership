package transport

import (
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/transport/http/authHandler"
	"orderPickupPoint/internal/transport/http/pickupPointHandler"
	"orderPickupPoint/internal/transport/http/receptionHandler"

	"github.com/gorilla/mux"
)

type Handler struct {
	Services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		Services: services,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()

	authHandler := authHandler.NewAuthHandler(h.Services.Auth)
	receptionHandler := receptionHandler.NewReceptionHandler(h.Services.Reception)
	pupHandler := pickupPointHandler.NewPickupPointHandler(h.Services.PickupPoint)

	modOnly := []string{"moderator"}
	router.Handle("/", authHandler.IsAvaliableRoleMiddleware(authHandler.IsSignedInMiddleware(receptionHandler.IsWorking), modOnly))

	router.HandleFunc("/dummyLogin", authHandler.DummyLogin).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// router.HandleFunc("/pvz", authHandler.IsSignedInMiddleware(authHandler.IsAvaliableRoleMiddleware(pupHandler.Create, modOnly))).Methods("POST")
	router.HandleFunc("/pvz", authHandler.IsAvaliableRoleMiddleware(authHandler.IsSignedInMiddleware(pupHandler.Create), modOnly)).Methods("POST")

	//router.HandleFunc("/receptions", ).Methods("POST")

	return router
}
