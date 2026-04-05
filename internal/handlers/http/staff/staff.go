package staff_handler

import (
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"
	"agnos-backend/internal/utility"
	"log"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc port.IStaffService
}

func New(svc port.IStaffService) port.IStaffHandler {
	return &handler{svc: svc}
}

func (h *handler) Create(c *gin.Context) {
	var req models.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("create staff handler: bind request error: %v", err)
		utility.BadRequest(c, err.Error())
		return
	}

	staff, err := h.svc.Create(req.Username, req.Password, req.Hospital)
	if err != nil {
		utility.InternalServerError(c, err.Error())
		return
	}

	utility.Success(c, models.CreateStaffResponse{
		StaffID:    staff.ID.String(),
		HospitalID: staff.HospitalID.String(),
	})
}

func (h *handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("login handler: bind request error: %v", err)
		utility.BadRequest(c, err.Error())
		return
	}

	token, err := h.svc.Login(req.Username, req.Password, req.Hospital)
	if err != nil {
		log.Printf("login handler: service error username=%s: %v", req.Username, err)
		utility.Unauthorized(c, err.Error())
		return
	}

	utility.Success(c, models.LoginResponse{
		Token: token,
	})
}
