package patient_handler

import (
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"
	"agnos-backend/internal/utility"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	svc port.IPatientService
}

func New(svc port.IPatientService) port.IPatientHandler {
	return &handler{svc: svc}
}

func (h *handler) Search(c *gin.Context) {
	hospitalIDStr, exists := c.Get("hospital_id")
	if !exists {
		log.Printf("search patient handler: hospital_id not found in token")
		utility.Unauthorized(c, "unauthorized")
		return
	}

	hospitalID, err := uuid.Parse(hospitalIDStr.(string))
	if err != nil {
		log.Printf("search patient handler: invalid hospital_id=%v: %v", hospitalIDStr, err)
		utility.BadRequest(c, "invalid hospital id in token")
		return
	}

	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("search patient handler: bind request error: %v", err)
		utility.BadRequest(c, err.Error())
		return
	}

	patients, err := h.svc.Search(hospitalID, req)
	if err != nil {
		utility.InternalServerError(c, err.Error())
		return
	}

	utility.Success(c, models.SearchResponse{
		Count:    len(patients),
		Patients: patients,
	})
}
