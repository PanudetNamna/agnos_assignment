package port

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=handlers.go -destination=./mocks/mock_handlers.go -package=mocks

type IPatientHandler interface {
	Search(c *gin.Context)
}

type IStaffHandler interface {
	Create(c *gin.Context)
	Login(c *gin.Context)
}
