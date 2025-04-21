package pickupPointHandler

import (
	"encoding/json"
	"net/http"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/utils/errorsHandl"
	"strconv"
	"time"
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
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var pickupPoint *models.PickupPointAPI
	if err := json.NewDecoder(r.Body).Decode(&pickupPoint); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	pickupPoint, err := h.pickupPointService.Create(r.Context(), pickupPoint)
	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pickupPoint); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
	}

}

func (h *PickupPointHandler) GetReceptionsInfo(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	page := r.URL.Query().Get("page")
	pageLimit := r.URL.Query().Get("limit")

	filter := &models.PvzFilter{}

	if startDate != "" && endDate != "" {
		if sd, err := time.Parse(time.RFC3339Nano, startDate); err == nil {
			filter.StartDate = &sd
		} else {
			errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
			return
		}
		if ed, err := time.Parse(time.RFC3339Nano, endDate); err == nil {
			filter.EndDate = &ed
		} else {
			errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
	if page != "" {
		val, err := strconv.Atoi(page)
		if err != nil || val < 1 {
			errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
			return
		}
		filter.Page = val
	} else {
		filter.Page = 1
	}

	if pageLimit != "" {
		val, err := strconv.Atoi(pageLimit)
		if err != nil || val < 1 || val > 30 {
			errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
			return
		}
		filter.PageLimit = val
	} else {
		filter.PageLimit = 10
	}

	info, err := h.pickupPointService.GetInfo(r.Context(), filter)

	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
