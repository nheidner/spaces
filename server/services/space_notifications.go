package services

import (
	"context"
	"net/http"
	"spaces-p/common"
	"spaces-p/errors"
	"spaces-p/models"
	localmemory "spaces-p/repositories/local_memory"
	"spaces-p/uuid"
	"time"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
)

type SpaceNotificationsService struct {
	logger          common.Logger
	cacheRepo       common.CacheRepository
	localMemoryRepo *localmemory.LocalMemoryRepo
}

func NewSpaceNotificationsService(logger common.Logger, cacheRepo common.CacheRepository, localMemoryRepo *localmemory.LocalMemoryRepo) *SpaceNotificationsService {
	return &SpaceNotificationsService{logger, cacheRepo, localMemoryRepo}
}

func (ss *SpaceNotificationsService) SpaceConnect(ctx context.Context, c *gin.Context, spaceId uuid.Uuid, authenticatedUser models.User) error {
	const op errors.Op = "services.SpaceNotificationsService.SpaceConnect"

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		return errors.E(op, err, http.StatusBadRequest)
	}
	// CHECK: will the closing status not always be StatusInternalError?
	defer conn.Close(websocket.StatusInternalError, "")
	err = ss.subscribe(ctx, conn, spaceId, authenticatedUser.ID)
	return errors.E(op, err)
}

func (ss *SpaceNotificationsService) subscribe(ctx context.Context, conn *websocket.Conn, spaceId uuid.Uuid, userId models.UserUid) error {
	const op errors.Op = "services.SpaceNotificationsService.subscribe"

	ctx = conn.CloseRead(ctx)

	session := ss.localMemoryRepo.AddSession(localmemory.NewSessionInput{
		SpaceId:         spaceId,
		UserId:          userId,
		NotificationsCh: make(chan models.Notification, localmemory.NotificationsBufferSize),
		CloseSlow: func() {
			conn.Close(websocket.StatusInternalError, "")
		},
	})
	defer ss.localMemoryRepo.DeleteSession(session.SpaceId, session.SessionId)

	for {
		select {
		case notification := <-session.NotificationsCh:
			err := writeWithTimeout(ctx, 5*time.Second, conn, notification)
			if err != nil {
				return errors.E(op, err)
			}
		case <-ctx.Done():
			return errors.E(op, ctx.Err())
		}
	}
}

func writeWithTimeout(ctx context.Context, timeout time.Duration, conn *websocket.Conn, notification []byte) error {
	const op errors.Op = "services.writeWithTimeout"

	ctx, cancel := context.WithTimeout(ctx, 5*timeout)
	defer cancel()

	err := conn.Write(ctx, websocket.MessageText, notification)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
