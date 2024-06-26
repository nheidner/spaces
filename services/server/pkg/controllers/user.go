package controllers

import (
	"net/http"
	"spaces-p/pkg/common"
	"spaces-p/pkg/errors"
	"spaces-p/pkg/services"
	"spaces-p/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
	logger      common.Logger
	authClient  common.AuthClient
}

func NewUserController(logger common.Logger, userService *services.UserService, authClient common.AuthClient) *UserController {
	return &UserController{userService, logger, authClient}
}

func (uc *UserController) CreateUserFromIdToken(c *gin.Context) {
	const op errors.Op = "controllers.UserController.CreateUser"
	var ctx = c.Request.Context()
	var body struct {
		IdToken string `json:"idToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	user, err := uc.userService.CreateUserFromIdToken(ctx, uc.authClient, body.IdToken)
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (uc *UserController) GetUser(c *gin.Context) {
	const op errors.Op = "controllers.UserController.GetUser"
	var ctx = c.Request.Context()

	userId := utils.GetUserUidFromPath(c)

	user, err := uc.userService.GetUser(ctx, userId)
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (uc *UserController) GetAuthedUser(c *gin.Context) {
	const op errors.Op = "controllers.UserController.GetAuthedUser"

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
