package services

import (
	"context"
	"net/http"
	"spaces-p/pkg/common"
	"spaces-p/pkg/errors"
	"spaces-p/pkg/models"
	localmemory "spaces-p/pkg/repositories/local_memory"
	"spaces-p/pkg/uuid"
	"time"
)

type ThreadService struct {
	logger          common.Logger
	cacheRepo       common.CacheRepository
	localMemoryRepo *localmemory.LocalMemoryRepo
}

func NewThreadService(logger common.Logger, cacheRepo common.CacheRepository, localMemoryRepo *localmemory.LocalMemoryRepo) *ThreadService {
	return &ThreadService{logger, cacheRepo, localMemoryRepo}
}

func (ts *ThreadService) CreateThread(ctx context.Context, spaceId, parentMessageId uuid.Uuid, authenticatedUserId models.UserUid) (uuid.Uuid, error) {
	const op errors.Op = "services.ThreadService.CreateThread"

	m, err := ts.cacheRepo.GetMessage(ctx, parentMessageId)
	switch {
	case errors.Is(err, common.ErrNotFound):
		return uuid.Nil, errors.E(op, err, http.StatusBadRequest)
	case err != nil:
		return uuid.Nil, errors.E(op, err, http.StatusInternalServerError)
	case m.ChildThreadId != uuid.Nil:
		err := errors.New("the parent message's child thread id has been already set")
		return uuid.Nil, errors.E(op, err, http.StatusBadRequest)
	}

	var createdAt = time.Now()
	thread, err := ts.cacheRepo.SetThread(ctx, spaceId, parentMessageId, createdAt)
	if err != nil {
		return uuid.Nil, errors.E(op, err, http.StatusInternalServerError)
	}

	ts.localMemoryRepo.PublishNewThread(spaceId, authenticatedUserId, *thread)

	return thread.ID, nil
}

func (ts *ThreadService) CreateTopLevelThread(ctx context.Context, spaceId uuid.Uuid, newTopLevelThreadFirstMessage models.NewTopLevelThreadFirstMessage) (uuid.Uuid, uuid.Uuid, error) {
	const op errors.Op = "services.ThreadService.CreateTopLevelThread"

	createdTopLevelThread, createdFirstMessage, err := ts.cacheRepo.SetTopLevelThread(ctx, spaceId, newTopLevelThreadFirstMessage)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.E(op, err, http.StatusInternalServerError)
	}

	ts.localMemoryRepo.PublishNewToplevelThread(spaceId, newTopLevelThreadFirstMessage.SenderId, *createdTopLevelThread)

	return createdTopLevelThread.ID, createdFirstMessage.ID, nil
}
