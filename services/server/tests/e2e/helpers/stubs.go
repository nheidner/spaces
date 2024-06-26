//go:build e2e
// +build e2e

package helpers

import (
	"context"
	"spaces-p/pkg/common"
	"spaces-p/pkg/models"
	"time"
)

// implements the common.AuthClient interface
type StubAuthClient struct {
	currentTestUserId models.UserUid
}

func (tac *StubAuthClient) VerifyToken(ctx context.Context, idToken string) (*common.UserTokenData, error) {
	return &common.UserTokenData{
		SignInProvider:  "email",
		EmailIsVerified: true,
		UserId:          tac.currentTestUserId,
	}, nil
}

func (tac *StubAuthClient) CreateUser(ctx context.Context, email, password string, emailIsVerified bool) (models.UserUid, error) {
	return "", nil
}

func (tac *StubAuthClient) DeleteAllUsers(ctx context.Context) (int, error) {
	return 0, nil
}

func (tac *StubAuthClient) SetCurrentTestUser(newCurrentTestUserId models.UserUid) {
	tac.currentTestUserId = newCurrentTestUserId
}

// implements the common.GeocodeRepository interface
type SpyGeocodeRepository struct {
	count              int
	currentTestAddress *models.Address
	currentErr         error
}

func (gr *SpyGeocodeRepository) GetAddress(ctx context.Context, location models.Location) (*models.Address, error) {
	gr.count++

	return gr.currentTestAddress, gr.currentErr
}

func (gr *SpyGeocodeRepository) SetTestAddress(newCurrentTestAddress models.Address, err error) {
	gr.currentTestAddress = &newCurrentTestAddress
	gr.currentErr = err
}

func (gr *SpyGeocodeRepository) Count() int {
	return gr.count
}

func (gr *SpyGeocodeRepository) Reset() {
	gr.count = 0
	gr.currentTestAddress = nil
	gr.currentErr = nil
}

// implements the common.Logger interface
type NoopLogger struct{}

func (lg *NoopLogger) Info(v ...any)  {}
func (lg *NoopLogger) Error(v ...any) {}
func (lg *NoopLogger) RequestInfo(method, path, clientIP string, statusCode int, latency time.Duration) {
}
