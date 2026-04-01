package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/service"
	orderservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeOrderService struct {
	createResp   *orderservicev1.CreateOrderResponse
	createErr    error
	getResp      *orderservicev1.GetOrderResponse
	getErr       error
	listResp     *orderservicev1.ListOrdersResponse
	listErr      error
	initiateResp *orderservicev1.InitiatePaymentResponse
	initiateErr  error
	confirmResp  *orderservicev1.ConfirmPaymentResponse
	confirmErr   error
	cancelResp   *orderservicev1.CancelOrderResponse
	cancelErr    error
	statusResp   *orderservicev1.GetOrderStatusResponse
	statusErr    error
}

func (f *fakeOrderService) CreateOrder(context.Context, *orderservicev1.CreateOrderRequest) (*orderservicev1.CreateOrderResponse, error) {
	return f.createResp, f.createErr
}
func (f *fakeOrderService) GetOrder(context.Context, *orderservicev1.GetOrderRequest) (*orderservicev1.GetOrderResponse, error) {
	return f.getResp, f.getErr
}
func (f *fakeOrderService) ListOrders(context.Context, *orderservicev1.ListOrdersRequest) (*orderservicev1.ListOrdersResponse, error) {
	return f.listResp, f.listErr
}
func (f *fakeOrderService) InitiatePayment(context.Context, *orderservicev1.InitiatePaymentRequest) (*orderservicev1.InitiatePaymentResponse, error) {
	return f.initiateResp, f.initiateErr
}
func (f *fakeOrderService) ConfirmPayment(context.Context, *orderservicev1.ConfirmPaymentRequest) (*orderservicev1.ConfirmPaymentResponse, error) {
	return f.confirmResp, f.confirmErr
}
func (f *fakeOrderService) CancelOrder(context.Context, *orderservicev1.CancelOrderRequest) (*orderservicev1.CancelOrderResponse, error) {
	return f.cancelResp, f.cancelErr
}
func (f *fakeOrderService) GetOrderStatus(context.Context, *orderservicev1.GetOrderStatusRequest) (*orderservicev1.GetOrderStatusResponse, error) {
	return f.statusResp, f.statusErr
}

func TestOrderHandlerDelegatesSuccess(t *testing.T) {
	handler := NewOrderHandler(&fakeOrderService{
		createResp: &orderservicev1.CreateOrderResponse{Message: "ok"},
		statusResp: &orderservicev1.GetOrderStatusResponse{OrderId: "order-1"},
	})

	createResp, err := handler.CreateOrder(context.Background(), &orderservicev1.CreateOrderRequest{})
	require.NoError(t, err)
	require.Equal(t, "ok", createResp.Message)

	statusResp, err := handler.GetOrderStatus(context.Background(), &orderservicev1.GetOrderStatusRequest{})
	require.NoError(t, err)
	require.Equal(t, "order-1", statusResp.OrderId)
}

func TestOrderHandlerMapsErrors(t *testing.T) {
	handler := NewOrderHandler(&fakeOrderService{createErr: service.ErrInvalidArgument})
	_, err := handler.CreateOrder(context.Background(), &orderservicev1.CreateOrderRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	require.Equal(t, codes.NotFound, status.Code(mapError(service.ErrNotFound)))
	require.Equal(t, codes.AlreadyExists, status.Code(mapError(service.ErrAlreadyExists)))
	require.Equal(t, codes.FailedPrecondition, status.Code(mapError(service.ErrInvalidTransition)))
	require.Equal(t, codes.Aborted, status.Code(mapError(service.ErrPaymentFailed)))
	require.Equal(t, codes.Unimplemented, status.Code(mapError(service.ErrNotImplemented)))
	require.Equal(t, codes.Internal, status.Code(mapError(errors.New("boom"))))
}
