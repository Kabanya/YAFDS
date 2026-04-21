package app

import (
	"net/http"

	"github.com/Kabanya/YAFDS/courier/internal/usecase"
	"github.com/Kabanya/YAFDS/pkg/utils"
)

func NewCouriersHandler(courUC usecase.CourierUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		couriers, err := courUC.ListCouriers(r.Context())
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, couriers, http.StatusOK)
	}
}
