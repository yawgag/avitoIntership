package receptionHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/utils/errorsHandl"
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

	var reqData *models.Reception
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		errorsHandl.SendJsonError(w, "Bad request", http.StatusBadRequest)
		return
	}
}
