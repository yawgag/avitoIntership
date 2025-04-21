package receptionHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/utils/errorsHandl"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReceptionHandler struct {
	receptionService service.Reception
}

func NewReceptionHandler(receptionService service.Reception) *ReceptionHandler {
	return &ReceptionHandler{
		receptionService: receptionService,
	}
}

func (h *ReceptionHandler) IsWorking(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "is working")
}

func (h *ReceptionHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var reception *models.ReceptionAPI
	if err := json.NewDecoder(r.Body).Decode(&reception); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	reception, err := h.receptionService.CreateReception(r.Context(), reception.PickupPointId)
	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reception); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
	}

}

func (h *ReceptionHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	var productAPI *models.ProductAPI
	if err := json.NewDecoder(r.Body).Decode(&productAPI); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	product, err := h.receptionService.AddProduct(r.Context(), productAPI)
	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
	}
}

func (h *ReceptionHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pvzId, err := uuid.Parse(vars["pvzId"])
	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.receptionService.DeleteLastProductInReception(r.Context(), pvzId)
	if err != nil {
		errorsHandl.SendJsonError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func (h *ReceptionHandler) CloseReception(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pvzId, err := uuid.Parse(vars["pvzId"])
	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.receptionService.CloseReception(r.Context(), pvzId)
	if err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}
}
