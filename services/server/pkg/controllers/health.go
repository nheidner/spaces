package controllers

import (
	"net/http"
	"spaces-p/pkg/common"
	"spaces-p/pkg/errors"
	"spaces-p/pkg/services"
	"spaces-p/pkg/utils"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	logger        common.Logger
	healthService *services.HealthService
}

func NewHealthController(logger common.Logger, healthService *services.HealthService) *HealthController {
	return &HealthController{logger, healthService}
}
func (hs *HealthController) HealthCheck(c *gin.Context) {
	const op errors.Op = "controllers.HealthController.HealthCheck"

	if err := hs.healthService.GetDbHealth(c); err != nil {
		utils.WriteError(c, errors.E(op, err), hs.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "db": "OK"})
}
