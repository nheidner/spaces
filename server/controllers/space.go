package controllers

import (
	"net/http"
	"spaces-p/common"
	"spaces-p/errors"
	"spaces-p/models"
	"spaces-p/services"
	"spaces-p/utils"

	"github.com/gin-gonic/gin"
)

type SpaceController struct {
	logger         common.Logger
	spaceService   *services.SpaceService
	threadService  *services.ThreadService
	messageService *services.MessageService
}

func NewSpaceController(logger common.Logger, spaceService *services.SpaceService, threadService *services.ThreadService, messageService *services.MessageService) *SpaceController {
	return &SpaceController{logger, spaceService, threadService, messageService}
}

func (uc *SpaceController) GetSpaces(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.GetSpaces"
	var ctx = c.Request.Context()
	var query struct {
		Location string  `form:"searchByLocation"`
		Radius   float64 `form:"radius"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	var location models.Location
	if err := location.ParseString(query.Location); err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	spaces, err := uc.spaceService.GetSpacesByLocation(ctx, location, models.Radius(query.Radius))
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": spaces})
}

func (uc *SpaceController) CreateSpace(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.CreateSpace"
	var ctx = c.Request.Context()

	var body models.NewSpace
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	spaceId, err := uc.spaceService.CreateSpace(ctx, body)
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"spaceId": spaceId,
	}})
}

func (uc *SpaceController) GetTopLevelThreads(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.GetTopLevelThreads"
	var ctx = c.Request.Context()
	var query struct {
		Offset int64          `form:"offset"`
		Count  int64          `form:"count"`
		Sort   models.Sorting `form:"sort"`
	}
	if query.Count == 0 {
		query.Count = 10
	}

	spaceId, err := utils.GetSpaceIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	topLevelThreads, err := uc.spaceService.GetTopLevelThreads(ctx, spaceId, query.Sort, query.Offset, query.Count)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusInternalServerError), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": topLevelThreads})
}

func (uc *SpaceController) GetThreadWithMessages(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.GetThreadWithMessages"
	var ctx = c.Request.Context()
	var query struct {
		MessagesOffset int64          `form:"messages_offset"`
		MessagesCount  int64          `form:"messages_count"`
		MessagesSort   models.Sorting `form:"messages_sort"`
	}
	if query.MessagesCount == 0 {
		query.MessagesCount = 10
	}

	spaceId, err := utils.GetSpaceIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	threadId, err := utils.GetThreadIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	threads, err := uc.spaceService.GetThreadWithMessages(ctx, spaceId, threadId, query.MessagesSort, query.MessagesOffset, query.MessagesCount)
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": threads})
}

func (uc *SpaceController) CreateTopLevelThread(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.CreateTopLevelThread"
	var ctx = c.Request.Context()

	var body models.NewMessageInput
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	authenticatedUser, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusInternalServerError), uc.logger)
		return
	}

	spaceId, err := utils.GetSpaceIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	threadId, err := uc.threadService.CreateTopLevelThread(ctx, spaceId, models.NewTopLevelThreadFirstMessage{
		NewMessageInput: body,
		SenderId:        authenticatedUser.ID,
	})
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": threadId})
}

func (uc *SpaceController) CreateThread(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.CreateThread"
	var ctx = c.Request.Context()

	spaceId, err := utils.GetSpaceIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	messageId, err := utils.GetMessageIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	threadId, err := uc.threadService.CreateThread(ctx, spaceId, messageId)
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": threadId})
}

func (uc *SpaceController) CreateMessage(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.CreateMessage"
	var ctx = c.Request.Context()

	var body models.NewMessageInput
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	authenticatedUser, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusInternalServerError), uc.logger)
		return
	}

	spaceId, err := utils.GetSpaceIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	threadId, err := utils.GetThreadIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	messageId, err := uc.messageService.CreateMessage(ctx, spaceId, models.NewMessage{
		BaseMessage: models.BaseMessage(body),
		SenderId:    authenticatedUser.ID,
		ThreadId:    threadId,
	})
	if err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": messageId})
}

func (uc *SpaceController) LikeMessage(c *gin.Context) {
	const op errors.Op = "controllers.SpaceController.LikeMessage"
	var ctx = c.Request.Context()

	messageId, err := utils.GetMessageIdFromPath(c)
	if err != nil {
		utils.WriteError(c, errors.E(op, err, http.StatusBadRequest), uc.logger)
		return
	}

	if err := uc.messageService.LikeMessage(ctx, messageId); err != nil {
		utils.WriteError(c, errors.E(op, err), uc.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
