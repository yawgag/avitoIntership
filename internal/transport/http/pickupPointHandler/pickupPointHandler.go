package pickupPointHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/utils/errorsHandl"
)

type PickupPointHandler struct {
	pickupPointService service.PickupPoint
}

func NewPickupPointHandler(pickupPointService service.PickupPoint) *PickupPointHandler {
	return &PickupPointHandler{
		pickupPointService: pickupPointService,
	}
}

func (h *PickupPointHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusUnauthorized)
		return
	}

	var pickupPoint *models.PickupPointAPI
	if err := json.NewDecoder(r.Body).Decode(&pickupPoint); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusUnauthorized)
		return
	}

	pickupPoint, err := h.pickupPointService.Create(r.Context(), pickupPoint)
	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pickupPoint); err != nil {
		fmt.Println("err: ", err)
		//TODO: add json error
	}
}
