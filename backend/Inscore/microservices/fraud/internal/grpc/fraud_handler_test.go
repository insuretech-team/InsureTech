package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/fraud/internal/service"
	fraudv1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/entity/v1"
	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handlerRuleRepo struct{}

func (handlerRuleRepo) Create(context.Context, *fraudv1.FraudRule) error { return nil }
func (handlerRuleRepo) GetByID(context.Context, string) (*fraudv1.FraudRule, error) {
	return nil, repository.ErrRuleNotFound
}
func (handlerRuleRepo) Update(context.Context, string, *fraudv1.FraudRule) error { return nil }
func (handlerRuleRepo) List(context.Context, fraudv1.RuleCategory, bool, int, int) ([]*fraudv1.FraudRule, int32, error) {
	return nil, 0, nil
}
func (handlerRuleRepo) SetActive(context.Context, string, bool) error { return nil }

type handlerAlertRepo struct {
	alert  *fraudv1.FraudAlert
	getErr error
}

func (r handlerAlertRepo) Create(context.Context, *fraudv1.FraudAlert) error { return nil }
func (r handlerAlertRepo) GetByID(context.Context, string) (*fraudv1.FraudAlert, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	return r.alert, nil
}
func (handlerAlertRepo) List(context.Context, string, string, *time.Time, *time.Time, int, int) ([]*fraudv1.FraudAlert, int32, error) {
	return nil, 0, nil
}
func (handlerAlertRepo) UpdateStatus(context.Context, string, fraudv1.AlertStatus, string) error {
	return nil
}

type handlerCaseRepo struct{}

func (handlerCaseRepo) Create(context.Context, *fraudv1.FraudCase) error { return nil }
func (handlerCaseRepo) GetByID(context.Context, string) (*fraudv1.FraudCase, error) {
	return nil, repository.ErrCaseNotFound
}
func (handlerCaseRepo) Update(context.Context, string, fraudv1.CaseStatus, fraudv1.CaseOutcome, string, string) error {
	return nil
}

func TestFraudHandlerCheckFraudMapsInvalidArgument(t *testing.T) {
	handler := NewFraudHandler(service.NewFraudService(handlerRuleRepo{}, handlerAlertRepo{}, handlerCaseRepo{}, nil))

	_, err := handler.CheckFraud(context.Background(), &fraudservicev1.CheckFraudRequest{})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestFraudHandlerGetFraudAlertSuccessAndNotFound(t *testing.T) {
	successHandler := NewFraudHandler(service.NewFraudService(
		handlerRuleRepo{},
		handlerAlertRepo{alert: &fraudv1.FraudAlert{Id: "alert-1"}},
		handlerCaseRepo{},
		nil,
	))

	resp, err := successHandler.GetFraudAlert(context.Background(), &fraudservicev1.GetFraudAlertRequest{
		FraudAlertId: "alert-1",
	})
	require.NoError(t, err)
	require.Equal(t, "alert-1", resp.FraudAlert.Id)

	notFoundHandler := NewFraudHandler(service.NewFraudService(
		handlerRuleRepo{},
		handlerAlertRepo{getErr: repository.ErrAlertNotFound},
		handlerCaseRepo{},
		nil,
	))

	_, err = notFoundHandler.GetFraudAlert(context.Background(), &fraudservicev1.GetFraudAlertRequest{
		FraudAlertId: "missing",
	})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestMapError(t *testing.T) {
	require.Equal(t, codes.InvalidArgument, status.Code(mapError(service.ErrInvalidArgument)))
	require.Equal(t, codes.NotFound, status.Code(mapError(service.ErrNotFound)))
	require.Equal(t, codes.Internal, status.Code(mapError(errors.New("boom"))))
}
