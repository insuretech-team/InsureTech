package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/events"
	paymentrepo "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/repository"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	paymententityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	paymenteventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/events/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentService struct {
	repo      domain.PaymentRepository
	publisher *events.Publisher
	config    *config.Config
	gateway   domain.PaymentGateway
}

func NewPaymentService(repo domain.PaymentRepository, publisher *events.Publisher, cfg *config.Config, gateway domain.PaymentGateway) *PaymentService {
	return &PaymentService{
		repo:      repo,
		publisher: publisher,
		config:    cfg,
		gateway:   gateway,
	}
}

func (s *PaymentService) InitiatePayment(ctx context.Context, req *paymentservicev1.InitiatePaymentRequest) (*paymentservicev1.InitiatePaymentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetIdempotencyKey()) == "" {
		return nil, fmt.Errorf("%w: idempotency_key is required", ErrInvalidArgument)
	}
	if req.GetAmount() == nil || req.GetAmount().GetAmount() <= 0 {
		return nil, fmt.Errorf("%w: amount is required", ErrInvalidArgument)
	}

	existing, err := s.repo.GetPaymentByIdempotencyKey(ctx, req.GetIdempotencyKey())
	if err == nil && existing != nil {
		return &paymentservicev1.InitiatePaymentResponse{
			PaymentId:      existing.GetPaymentId(),
			TransactionId:  existing.GetTransactionId(),
			PaymentUrl:     firstNonEmpty(gatewayValue(existing.GetGatewayResponse(), "payment_url"), gatewayValue(existing.GetGatewayResponse(), "GatewayPageURL")),
			Status:         existing.GetStatus().String(),
			ExpiresAt:      existing.GetInitiatedAt(),
			Provider:       existing.GetProvider(),
			GatewayPageUrl: gatewayValue(existing.GetGatewayResponse(), "GatewayPageURL"),
			TranId:         firstNonEmpty(existing.GetTranId(), existing.GetTransactionId()),
			SessionKey:     existing.GetSessionKey(),
		}, nil
	}
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}

	userID := resolveUserID(ctx, req.GetUserId())
	if userID == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrInvalidArgument)
	}

	method, gatewayName, err := normalizeMethod(req.GetPaymentMethod())
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(30 * time.Minute)
	paymentID := uuid.NewString()
	transactionID := merchantTransactionID(paymentID)
	currency := firstNonEmpty(req.GetCurrency(), req.GetAmount().GetCurrency(), "BDT")
	amount := cloneMoney(req.GetAmount(), currency)
	status := paymententityv1.PaymentStatus_PAYMENT_STATUS_PENDING
	paymentURL := ""
	if gatewayName == "SSLCOMMERZ" {
		status = paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED
	}

	// Resolve typed linkage fields — prefer new typed fields, fall back to metadata for compat
	orderID := firstNonEmpty(req.GetOrderId(), strings.TrimSpace(req.GetMetadata()["order_id"]))
	tenantID := firstNonEmpty(req.GetTenantId(), strings.TrimSpace(req.GetMetadata()["tenant_id"]))
	customerID := firstNonEmpty(req.GetCustomerId(), strings.TrimSpace(req.GetMetadata()["customer_id"]))
	organisationID := firstNonEmpty(req.GetOrganisationId(), strings.TrimSpace(req.GetMetadata()["organisation_id"]))
	invoiceID := firstNonEmpty(req.GetInvoiceId(), strings.TrimSpace(req.GetMetadata()["invoice_id"]))
	purchaseOrderID := firstNonEmpty(req.GetPurchaseOrderId(), strings.TrimSpace(req.GetMetadata()["purchase_order_id"]))

	gatewayPayload := map[string]string{
		"provider":       gatewayName,
		"callback_url":   req.GetCallbackUrl(),
		"order_id":       orderID,
		"tenant_id":      tenantID,
		"transaction_id": transactionID,
	}

	payment := &paymententityv1.Payment{
		PaymentId:       paymentID,
		TransactionId:   transactionID,
		PolicyId:        req.GetPolicyId(),
		Type:            paymententityv1.PaymentType_PAYMENT_TYPE_PREMIUM,
		Method:          method,
		Status:          status,
		Amount:          amount,
		Currency:        currency,
		PayerId:         userID,
		PayeeId:         normalizeOptionalUUID(s.config.DefaultPayeeID),
		InitiatedAt:     timestamppb.New(now),
		CreatedAt:       timestamppb.New(now),
		UpdatedAt:       timestamppb.New(now),
		Gateway:         gatewayName,
		Provider:        gatewayName,
		GatewayResponse: paymentrepo.MarshalGatewayResponse(gatewayPayload),
		IdempotencyKey:  req.GetIdempotencyKey(),
		// Typed linkage fields
		OrderId:            normalizeOptionalUUID(orderID),
		InvoiceId:          normalizeOptionalUUID(invoiceID),
		TenantId:           normalizeOptionalUUID(tenantID),
		CustomerId:         normalizeOptionalUUID(customerID),
		OrganisationId:     normalizeOptionalUUID(organisationID),
		PurchaseOrderId:    normalizeOptionalUUID(purchaseOrderID),
		ManualReviewStatus: paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_NOT_REQUIRED,
		ProviderReference:  transactionID,
		TranId:             transactionID,
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	if gatewayName == "SSLCOMMERZ" {
		if s.gateway == nil {
			payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_FAILED
			payment.FailureReason = "sslcommerz gateway client is not configured"
			payment.UpdatedAt = timestamppb.New(time.Now().UTC())
			_ = s.repo.UpdatePayment(ctx, payment)
			s.publishFailed(ctx, payment, "GATEWAY_NOT_CONFIGURED", payment.FailureReason)
			return nil, fmt.Errorf("%w: sslcommerz gateway client is not configured", ErrPaymentFailed)
		}

		successURL, failURL, cancelURL, ipnURL := s.buildGatewayURLs(paymentID, req.GetCallbackUrl())

		// Resolve customer details — prefer typed proto fields over metadata fallback
		sessionResp, err := s.gateway.InitSession(ctx, &domain.GatewaySessionRequest{
			PaymentID:        paymentID,
			TransactionID:    transactionID,
			Amount:           amount,
			Currency:         currency,
			SuccessURL:       successURL,
			FailURL:          failURL,
			CancelURL:        cancelURL,
			IPNURL:           ipnURL,
			OrderID:          orderID,
			TenantID:         tenantID,
			CustomerName:     firstNonEmpty(req.GetCustomerName(), strings.TrimSpace(req.GetMetadata()["customer_name"])),
			CustomerEmail:    firstNonEmpty(req.GetCustomerEmail(), strings.TrimSpace(req.GetMetadata()["customer_email"])),
			CustomerPhone:    firstNonEmpty(req.GetCustomerPhone(), strings.TrimSpace(req.GetMetadata()["customer_phone"])),
			CustomerAddr1:    firstNonEmpty(req.GetCustomerAddressLine1(), strings.TrimSpace(req.GetMetadata()["customer_address"])),
			CustomerCity:     firstNonEmpty(req.GetCustomerCity(), strings.TrimSpace(req.GetMetadata()["customer_city"])),
			CustomerPostcode: firstNonEmpty(req.GetCustomerPostcode(), strings.TrimSpace(req.GetMetadata()["customer_postcode"])),
			CustomerCountry:  firstNonEmpty(req.GetCustomerCountry(), strings.TrimSpace(req.GetMetadata()["customer_country"]), "Bangladesh"),
			Metadata:         req.GetMetadata(),
		})
		if err != nil {
			payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_FAILED
			payment.FailureReason = err.Error()
			payment.UpdatedAt = timestamppb.New(time.Now().UTC())
			payment.GatewayResponse = mergeGatewayResponse(payment.GatewayResponse, map[string]string{
				"provider_error": err.Error(),
			})
			_ = s.repo.UpdatePayment(ctx, payment)
			s.publishFailed(ctx, payment, "GATEWAY_INIT_FAILED", err.Error())
			return nil, fmt.Errorf("%w: %v", ErrPaymentFailed, err)
		}

		paymentURL = sessionResp.GatewayPageURL
		payment.SessionKey = sessionResp.SessionKey
		payment.ProviderReference = firstNonEmpty(sessionResp.SessionKey, sessionResp.TranID, payment.ProviderReference)
		payment.TranId = firstNonEmpty(sessionResp.TranID, payment.TranId)
		payment.GatewayResponse = mergeGatewayResponse(payment.GatewayResponse, sessionResp.RawFields)
		payment.GatewayResponse = mergeGatewayResponse(payment.GatewayResponse, map[string]string{
			"payment_url": paymentURL,
			"session_key": sessionResp.SessionKey,
			"tran_id":     sessionResp.TranID,
			"success_url": successURL,
			"fail_url":    failURL,
			"cancel_url":  cancelURL,
			"ipn_url":     ipnURL,
		})
		payment.UpdatedAt = timestamppb.New(time.Now().UTC())
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
	}

	refType, refID := referenceInfo(req)
	s.publishInitiated(ctx, payment, refType, refID)

	return &paymentservicev1.InitiatePaymentResponse{
		PaymentId:      paymentID,
		TransactionId:  transactionID,
		PaymentUrl:     paymentURL,
		Status:         status.String(),
		ExpiresAt:      timestamppb.New(expiresAt),
		Provider:       payment.GetProvider(),
		GatewayPageUrl: paymentURL,
		TranId:         payment.GetTranId(),
		SessionKey:     payment.GetSessionKey(),
	}, nil
}

func (s *PaymentService) VerifyPayment(ctx context.Context, req *paymentservicev1.VerifyPaymentRequest) (*paymentservicev1.VerifyPaymentResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}

	payment, err := s.repo.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, mapRepoError(err)
	}

	if payment.GetStatus() == paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS && !req.GetForceProviderRequery() {
		return &paymentservicev1.VerifyPaymentResponse{
			PaymentId: payment.GetPaymentId(),
			Status:    payment.GetStatus().String(),
			Payment:   payment,
			Verified:  true,
		}, nil
	}
	if isTerminalFailure(payment.GetStatus()) {
		return &paymentservicev1.VerifyPaymentResponse{
			PaymentId: payment.GetPaymentId(),
			Status:    payment.GetStatus().String(),
			Payment:   payment,
			Verified:  false,
			Error:     newError("PAYMENT_TERMINAL_STATE", "payment is already in terminal state"),
		}, nil
	}

	now := time.Now().UTC()
	if strings.EqualFold(firstNonEmpty(req.GetProvider(), payment.GetProvider(), payment.GetGateway()), "SSLCOMMERZ") {
		if s.gateway == nil {
			return nil, fmt.Errorf("%w: sslcommerz gateway client is not configured", ErrPaymentFailed)
		}

		sessionKey := firstNonEmpty(req.GetSessionKey(), payment.GetSessionKey(), gatewayValue(payment.GetGatewayResponse(), "session_key"), gatewayValue(payment.GetGatewayResponse(), "sessionkey"))
		validationReq := &domain.GatewayValidationRequest{
			PaymentID:     payment.GetPaymentId(),
			TransactionID: firstNonEmpty(req.GetValId(), req.GetTranId(), req.GetTransactionId(), payment.GetValId(), payment.GetTranId(), payment.GetTransactionId()),
			SessionKey:    sessionKey,
		}

		validationResp, err := s.gateway.ValidatePayment(ctx, validationReq)
		if err != nil || !providerStatusValid(validationResp) {
			validationResp, err = s.gateway.QueryPayment(ctx, &domain.GatewayValidationRequest{
				PaymentID:     payment.GetPaymentId(),
				TransactionID: payment.GetTransactionId(),
				SessionKey:    sessionKey,
			})
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrPaymentFailed, err)
			}
		}

		if !providerStatusValid(validationResp) {
			return &paymentservicev1.VerifyPaymentResponse{
				PaymentId: payment.GetPaymentId(),
				Status:    payment.GetStatus().String(),
				Payment:   payment,
				Verified:  false,
				Error:     newError("PAYMENT_NOT_VALIDATED", "provider validation did not confirm settlement"),
			}, nil
		}

		if !validationMatches(payment, validationResp) {
			return &paymentservicev1.VerifyPaymentResponse{
				PaymentId: payment.GetPaymentId(),
				Status:    payment.GetStatus().String(),
				Payment:   payment,
				Verified:  false,
				Error:     newError("PAYMENT_MISMATCH", "provider response does not match payment intent"),
			}, nil
		}

		payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS
		payment.TransactionId = firstNonEmpty(validationResp.TransactionID, req.GetTransactionId(), payment.GetTransactionId())
		payment.Provider = "SSLCOMMERZ"
		payment.ProviderReference = firstNonEmpty(validationResp.ValidationID, validationResp.TransactionID, sessionKey, payment.GetProviderReference())
		payment.TranId = firstNonEmpty(validationResp.TransactionID, req.GetTranId(), payment.GetTranId(), payment.GetTransactionId())
		payment.ValId = firstNonEmpty(validationResp.ValidationID, req.GetValId(), payment.GetValId())
		payment.SessionKey = firstNonEmpty(sessionKey, payment.GetSessionKey())
		payment.BankTranId = firstNonEmpty(validationResp.BankTransactionID, payment.GetBankTranId())
		payment.CardType = firstNonEmpty(validationResp.CardType, payment.GetCardType())
		payment.CardBrand = firstNonEmpty(validationResp.CardBrand, payment.GetCardBrand())
		payment.CardIssuer = firstNonEmpty(validationResp.CardIssuer, payment.GetCardIssuer())
		payment.CardIssuerCountry = firstNonEmpty(validationResp.CardIssuerCountry, payment.GetCardIssuerCountry())
		validatedAt := validationResp.ValidatedAt
		if validatedAt.IsZero() {
			validatedAt = now
		}
		payment.ValidatedAt = timestamppb.New(validatedAt)
		payment.ValidationStatus = firstNonEmpty(validationResp.Status, payment.GetValidationStatus())
		payment.RiskLevel = firstNonEmpty(validationResp.RiskLevel, payment.GetRiskLevel())
		payment.RiskTitle = firstNonEmpty(validationResp.RiskTitle, payment.GetRiskTitle())
		payment.CompletedAt = timestamppb.New(now)
		payment.UpdatedAt = timestamppb.New(now)
		payment.ReceiptUrl = s.receiptURL(payment.GetPaymentId())
		payment.GatewayResponse = mergeGatewayResponse(payment.GetGatewayResponse(), validationResp.RawFields)
		payment.GatewayResponse = mergeGatewayResponse(payment.GetGatewayResponse(), map[string]string{
			"verified_at":    now.Format(time.RFC3339Nano),
			"validation_id":  validationResp.ValidationID,
			"val_id":         validationResp.ValidationID,
			"bank_tran_id":   validationResp.BankTransactionID,
			"risk_level":     validationResp.RiskLevel,
			"risk_title":     validationResp.RiskTitle,
			"card_type":      validationResp.CardType,
			"card_brand":     validationResp.CardBrand,
			"card_issuer":    validationResp.CardIssuer,
			"issuer_country": validationResp.CardIssuerCountry,
		})

		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
		s.publishCompleted(ctx, payment, firstNonEmpty(req.GetIdempotencyKey(), payment.GetIdempotencyKey()))
		return &paymentservicev1.VerifyPaymentResponse{
			PaymentId: payment.GetPaymentId(),
			Status:    payment.GetStatus().String(),
			Payment:   payment,
			Verified:  true,
		}, nil
	}

	if transactionID := strings.TrimSpace(req.GetTransactionId()); transactionID != "" {
		payment.TransactionId = transactionID
	}
	payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS
	payment.CompletedAt = timestamppb.New(now)
	payment.UpdatedAt = timestamppb.New(now)
	payment.ValidatedAt = timestamppb.New(now)
	payment.ValidationStatus = "MANUAL"
	payment.GatewayResponse = mergeGatewayResponse(payment.GetGatewayResponse(), map[string]string{
		"verification_method": strings.ToUpper(req.GetPaymentMethod()),
		"verified_at":         now.Format(time.RFC3339Nano),
	})
	payment.ReceiptUrl = s.receiptURL(payment.GetPaymentId())

	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}
	s.publishCompleted(ctx, payment, firstNonEmpty(req.GetIdempotencyKey(), payment.GetIdempotencyKey()))

	return &paymentservicev1.VerifyPaymentResponse{
		PaymentId: payment.GetPaymentId(),
		Status:    payment.GetStatus().String(),
		Payment:   payment,
		Verified:  true,
	}, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, req *paymentservicev1.GetPaymentRequest) (*paymentservicev1.GetPaymentResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}
	payment, err := s.repo.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, mapRepoError(err)
	}
	return &paymentservicev1.GetPaymentResponse{Payment: payment}, nil
}

func (s *PaymentService) ListPayments(ctx context.Context, req *paymentservicev1.ListPaymentsRequest) (*paymentservicev1.ListPaymentsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}

	offset, err := parsePageToken(req.GetPageToken())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid page_token", ErrInvalidArgument)
	}

	filters := domain.PaymentFilters{
		UserID:    strings.TrimSpace(req.GetUserId()),
		PolicyID:  strings.TrimSpace(req.GetPolicyId()),
		Status:    strings.TrimSpace(req.GetStatus()),
		Limit:     normalizePageSize(req.GetPageSize()),
		Offset:    offset,
		StartDate: tsPtr(req.GetStartDate()),
		EndDate:   tsPtr(req.GetEndDate()),
	}

	payments, total, err := s.repo.ListPayments(ctx, filters)
	if err != nil {
		return nil, err
	}

	nextPageToken := ""
	if int64(offset)+int64(len(payments)) < total {
		nextPageToken = fmt.Sprintf("%d", offset+len(payments))
	}

	return &paymentservicev1.ListPaymentsResponse{
		Payments:      payments,
		NextPageToken: nextPageToken,
		TotalCount:    int32(total),
	}, nil
}

func (s *PaymentService) InitiateRefund(ctx context.Context, req *paymentservicev1.InitiateRefundRequest) (*paymentservicev1.InitiateRefundResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}

	payment, err := s.repo.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, mapRepoError(err)
	}
	if payment.GetStatus() != paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS {
		return nil, fmt.Errorf("%w: only successful payments can be refunded", ErrInvalidTransition)
	}

	refundAmount := req.GetRefundAmount()
	if refundAmount == nil {
		refundAmount = cloneMoney(payment.GetAmount(), payment.GetCurrency())
	}

	now := time.Now().UTC()
	refundStatus := paymententityv1.PaymentRefundStatus_PAYMENT_REFUND_STATUS_COMPLETED
	if payment.GetGateway() == "SSLCOMMERZ" {
		if s.gateway == nil {
			return nil, fmt.Errorf("%w: sslcommerz gateway client is not configured", ErrPaymentFailed)
		}
		bankTranID := firstNonEmpty(gatewayValue(payment.GetGatewayResponse(), "bank_tran_id"), gatewayValue(payment.GetGatewayResponse(), "bankTranId"))
		if strings.TrimSpace(bankTranID) == "" {
			return nil, fmt.Errorf("%w: bank_tran_id is required for sslcommerz refund", ErrInvalidArgument)
		}
		refundResp, err := s.gateway.InitiateRefund(ctx, &domain.GatewayRefundRequest{
			PaymentID:         payment.GetPaymentId(),
			BankTransactionID: bankTranID,
			Amount:            cloneMoney(refundAmount, payment.GetCurrency()),
			Reason:            req.GetReason(),
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPaymentFailed, err)
		}
		refundStatus = paymententityv1.PaymentRefundStatus_PAYMENT_REFUND_STATUS_PROCESSING
		if strings.Contains(strings.ToUpper(refundResp.Status), "REFUND") {
			refundStatus = paymententityv1.PaymentRefundStatus_PAYMENT_REFUND_STATUS_COMPLETED
		}
	}

	refund := &paymententityv1.PaymentRefund{
		RefundId:        uuid.NewString(),
		PaymentId:       payment.GetPaymentId(),
		RefundPaymentId: "",
		RefundAmount:    cloneMoney(refundAmount, payment.GetCurrency()),
		Reason:          req.GetReason(),
		Status:          refundStatus,
		ApprovedBy:      firstNonEmpty(req.GetInitiatedBy(), resolveUserID(ctx, "")),
		ApprovedAt:      timestamppb.New(now),
		CreatedAt:       timestamppb.New(now),
		UpdatedAt:       timestamppb.New(now),
	}
	if refundStatus == paymententityv1.PaymentRefundStatus_PAYMENT_REFUND_STATUS_COMPLETED {
		refund.ProcessedAt = timestamppb.New(now)
	}

	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		return nil, err
	}

	if refundStatus == paymententityv1.PaymentRefundStatus_PAYMENT_REFUND_STATUS_COMPLETED {
		payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_REFUNDED
		payment.UpdatedAt = timestamppb.New(now)
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
		s.publishRefundProcessed(ctx, refund, payment)
	}

	return &paymentservicev1.InitiateRefundResponse{
		RefundId:    refund.GetRefundId(),
		Status:      refund.GetStatus().String(),
		InitiatedAt: timestamppb.New(now),
	}, nil
}

func (s *PaymentService) GetRefundStatus(ctx context.Context, req *paymentservicev1.GetRefundStatusRequest) (*paymentservicev1.GetRefundStatusResponse, error) {
	if req == nil || strings.TrimSpace(req.GetRefundId()) == "" {
		return nil, fmt.Errorf("%w: refund_id is required", ErrInvalidArgument)
	}

	refund, err := s.repo.GetRefund(ctx, req.GetRefundId())
	if err != nil {
		return nil, mapRepoError(err)
	}

	return &paymentservicev1.GetRefundStatusResponse{
		RefundId:     refund.GetRefundId(),
		Status:       refund.GetStatus().String(),
		RefundAmount: refund.GetRefundAmount(),
		CompletedAt:  refund.GetProcessedAt(),
	}, nil
}

func (s *PaymentService) ListPaymentMethods(context.Context, *paymentservicev1.ListPaymentMethodsRequest) (*paymentservicev1.ListPaymentMethodsResponse, error) {
	methods := []*paymentservicev1.PaymentMethod{
		{
			MethodId:    "bank_transfer_b2b",
			MethodType:  "BANK_TRANSFER",
			DisplayName: "Bank Transfer",
			IsDefault:   false,
			Metadata:    map[string]string{"channel": "B2B", "provider": "BANK_TRANSFER"},
		},
		{
			MethodId:    "sslcommerz_card_b2c",
			MethodType:  "CARD",
			DisplayName: "SSLCommerz Card Checkout",
			IsDefault:   true,
			Metadata:    map[string]string{"channel": "B2C", "provider": "SSLCOMMERZ"},
		},
	}
	return &paymentservicev1.ListPaymentMethodsResponse{
		PaymentMethods: methods,
		TotalCount:     int32(len(methods)),
	}, nil
}

func (s *PaymentService) AddPaymentMethod(ctx context.Context, req *paymentservicev1.AddPaymentMethodRequest) (*paymentservicev1.AddPaymentMethodResponse, error) {
	if req == nil || req.GetPaymentMethod() == nil {
		return nil, fmt.Errorf("%w: payment_method is required", ErrInvalidArgument)
	}
	if _, _, err := normalizeMethod(req.GetPaymentMethod().GetMethodType()); err != nil {
		return nil, err
	}

	methodID := req.GetPaymentMethod().GetMethodId()
	if strings.TrimSpace(methodID) == "" {
		methodID = "pm_" + uuid.NewString()
	}

	return &paymentservicev1.AddPaymentMethodResponse{MethodId: methodID, Success: true}, nil
}

func (s *PaymentService) ReconcilePayments(ctx context.Context, req *paymentservicev1.ReconcilePaymentsRequest) (*paymentservicev1.ReconcilePaymentsResponse, error) {
	filters := domain.PaymentFilters{
		PaymentMethod: strings.TrimSpace(req.GetPaymentMethod()),
		StartDate:     tsPtr(req.GetStartDate()),
		EndDate:       tsPtr(req.GetEndDate()),
		Limit:         1000,
	}

	payments, total, err := s.repo.ListPayments(ctx, filters)
	if err != nil {
		return nil, err
	}

	var reconciled int32
	mismatchIDs := make([]string, 0)
	for _, payment := range payments {
		switch payment.GetStatus() {
		case paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS, paymententityv1.PaymentStatus_PAYMENT_STATUS_REFUNDED:
			reconciled++
		default:
			mismatchIDs = append(mismatchIDs, payment.GetPaymentId())
		}
	}

	return &paymentservicev1.ReconcilePaymentsResponse{
		TotalPayments:      int32(total),
		ReconciledCount:    reconciled,
		MismatchCount:      int32(len(mismatchIDs)),
		MismatchPaymentIds: mismatchIDs,
	}, nil
}

func (s *PaymentService) HandleGatewayWebhook(ctx context.Context, req *paymentservicev1.HandleGatewayWebhookRequest) (*paymentservicev1.HandleGatewayWebhookResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	provider := strings.ToUpper(strings.TrimSpace(req.GetProvider()))
	if provider == "" {
		return nil, fmt.Errorf("%w: provider is required", ErrInvalidArgument)
	}
	if provider != "SSLCOMMERZ" {
		return nil, fmt.Errorf("%w: provider %s is not supported", ErrNotImplemented, provider)
	}

	values, err := url.ParseQuery(string(req.GetRawPayload()))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid webhook payload", ErrInvalidArgument)
	}

	payment, err := s.lookupPaymentForCallback(ctx, provider, values)
	if err != nil {
		return nil, err
	}

	receivedAt := time.Now().UTC()
	if req.GetReceivedAt() != nil {
		receivedAt = req.GetReceivedAt().AsTime().UTC()
	}
	callbackType := strings.ToLower(strings.TrimSpace(req.GetHeaders()["x-payment-callback-type"]))
	if callbackType == "" {
		callbackType = "webhook"
	}

	payment.Provider = provider
	payment.ProviderReference = firstNonEmpty(values.Get("val_id"), values.Get("sessionkey"), values.Get("tran_id"), payment.GetProviderReference())
	payment.TranId = firstNonEmpty(values.Get("tran_id"), payment.GetTranId(), payment.GetTransactionId())
	payment.ValId = firstNonEmpty(values.Get("val_id"), payment.GetValId())
	payment.SessionKey = firstNonEmpty(values.Get("sessionkey"), payment.GetSessionKey())
	payment.BankTranId = firstNonEmpty(values.Get("bank_tran_id"), payment.GetBankTranId())
	payment.UpdatedAt = timestamppb.New(receivedAt)
	if callbackType == "webhook" || callbackType == "ipn" {
		payment.IpnReceivedAt = timestamppb.New(receivedAt)
	} else {
		payment.CallbackReceivedAt = timestamppb.New(receivedAt)
	}
	payment.GatewayResponse = mergeGatewayResponse(payment.GetGatewayResponse(), callbackFields(values))
	payment.GatewayResponse = mergeGatewayResponse(payment.GetGatewayResponse(), map[string]string{
		"last_callback_type": callbackType,
		"last_remote_addr":   req.GetRemoteAddr(),
	})

	switch callbackType {
	case "fail":
		if !isSettledStatus(payment.GetStatus()) {
			payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_FAILED
			payment.FailureReason = firstNonEmpty(values.Get("error"), values.Get("failedreason"), values.Get("status_message"), "sslcommerz failure callback received")
		}
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
		return &paymentservicev1.HandleGatewayWebhookResponse{
			Accepted:  true,
			PaymentId: payment.GetPaymentId(),
			Status:    payment.GetStatus().String(),
		}, nil
	case "cancel":
		if !isSettledStatus(payment.GetStatus()) {
			payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_CANCELLED
			payment.FailureReason = "sslcommerz cancel callback received"
		}
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
		return &paymentservicev1.HandleGatewayWebhookResponse{
			Accepted:  true,
			PaymentId: payment.GetPaymentId(),
			Status:    payment.GetStatus().String(),
		}, nil
	default:
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
		verifyResp, err := s.VerifyPayment(ctx, &paymentservicev1.VerifyPaymentRequest{
			PaymentId:            payment.GetPaymentId(),
			TransactionId:        firstNonEmpty(values.Get("val_id"), values.Get("tran_id"), payment.GetTransactionId()),
			PaymentMethod:        payment.GetMethod().String(),
			IdempotencyKey:       fmt.Sprintf("%s:%s:%s", strings.ToLower(provider), callbackType, firstNonEmpty(values.Get("val_id"), values.Get("tran_id"), payment.GetPaymentId())),
			Provider:             provider,
			ValId:                values.Get("val_id"),
			TranId:               values.Get("tran_id"),
			SessionKey:           firstNonEmpty(values.Get("sessionkey"), payment.GetSessionKey()),
			ForceProviderRequery: callbackType == "webhook" || callbackType == "ipn",
		})
		if err != nil {
			return nil, err
		}
		return &paymentservicev1.HandleGatewayWebhookResponse{
			Accepted:  verifyResp.GetVerified(),
			PaymentId: verifyResp.GetPaymentId(),
			Status:    verifyResp.GetStatus(),
		}, nil
	}
}

func (s *PaymentService) GetPaymentByProviderReference(ctx context.Context, req *paymentservicev1.GetPaymentByProviderReferenceRequest) (*paymentservicev1.GetPaymentByProviderReferenceResponse, error) {
	if req == nil || strings.TrimSpace(req.GetProvider()) == "" || strings.TrimSpace(req.GetProviderReference()) == "" {
		return nil, fmt.Errorf("%w: provider and provider_reference are required", ErrInvalidArgument)
	}
	payment, err := s.repo.GetPaymentByProviderReference(ctx, strings.ToUpper(strings.TrimSpace(req.GetProvider())), strings.TrimSpace(req.GetProviderReference()))
	if err != nil {
		return nil, mapRepoError(err)
	}
	return &paymentservicev1.GetPaymentByProviderReferenceResponse{Payment: payment}, nil
}

func (s *PaymentService) SubmitManualPaymentProof(ctx context.Context, req *paymentservicev1.SubmitManualPaymentProofRequest) (*paymentservicev1.SubmitManualPaymentProofResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}
	proofFileID := normalizeOptionalUUID(req.GetManualProofFileId())
	if proofFileID == "" {
		return nil, fmt.Errorf("%w: manual_proof_file_id must be a UUID", ErrInvalidArgument)
	}
	payment, err := s.repo.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, mapRepoError(err)
	}
	if payment.GetMethod() != paymententityv1.PaymentMethod_PAYMENT_METHOD_BANK_TRANSFER {
		return nil, fmt.Errorf("%w: manual proof is only supported for bank transfer payments", ErrInvalidTransition)
	}

	now := time.Now().UTC()
	payment.ManualProofFileId = proofFileID
	payment.ManualReviewStatus = paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_PENDING
	payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_MANUAL_REVIEW_REQUIRED
	payment.UpdatedAt = timestamppb.New(now)
	payment.GatewayResponse = mergeGatewayResponse(payment.GetGatewayResponse(), map[string]string{
		"manual_proof_submitted_by": firstNonEmpty(req.GetSubmittedBy(), resolveUserID(ctx, "")),
		"manual_proof_notes":        strings.TrimSpace(req.GetNotes()),
	})
	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}
	// Publish proof-submitted event so back-office queues and compliance can react.
	submittedBy := firstNonEmpty(req.GetSubmittedBy(), resolveUserID(ctx, ""))
	s.publishManualProofSubmitted(ctx, payment, submittedBy)

	return &paymentservicev1.SubmitManualPaymentProofResponse{
		PaymentId: payment.GetPaymentId(),
		Status:    payment.GetStatus().String(),
	}, nil
}

func (s *PaymentService) ReviewManualPayment(ctx context.Context, req *paymentservicev1.ReviewManualPaymentRequest) (*paymentservicev1.ReviewManualPaymentResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}
	reviewedBy := normalizeOptionalUUID(firstNonEmpty(req.GetReviewedBy(), resolveUserID(ctx, "")))
	if reviewedBy == "" {
		return nil, fmt.Errorf("%w: reviewed_by is required", ErrInvalidArgument)
	}
	payment, err := s.repo.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, mapRepoError(err)
	}
	if payment.GetManualReviewStatus() != paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_PENDING {
		return nil, fmt.Errorf("%w: payment is not awaiting manual review", ErrInvalidTransition)
	}

	now := time.Now().UTC()
	payment.VerifiedBy = reviewedBy
	payment.VerifiedAt = timestamppb.New(now)
	payment.UpdatedAt = timestamppb.New(now)
	payment.GatewayResponse = mergeGatewayResponse(payment.GetGatewayResponse(), map[string]string{
		"manual_review_notes": strings.TrimSpace(req.GetReviewNotes()),
	})
	if req.GetApproved() {
		payment.ManualReviewStatus = paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_APPROVED
		payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS
		payment.CompletedAt = timestamppb.New(now)
		payment.ValidatedAt = timestamppb.New(now)
		payment.ValidationStatus = "MANUAL_APPROVED"
		payment.RejectionReason = ""
		payment.FailureReason = ""
		payment.ReceiptUrl = s.receiptURL(payment.GetPaymentId())
		s.ensureReceiptFields(payment, now)
	} else {
		if strings.TrimSpace(req.GetRejectionReason()) == "" {
			return nil, fmt.Errorf("%w: rejection_reason is required when approved=false", ErrInvalidArgument)
		}
		payment.ManualReviewStatus = paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_REJECTED
		payment.Status = paymententityv1.PaymentStatus_PAYMENT_STATUS_FAILED
		payment.RejectionReason = strings.TrimSpace(req.GetRejectionReason())
		payment.FailureReason = payment.RejectionReason
	}
	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}
	return &paymentservicev1.ReviewManualPaymentResponse{
		PaymentId: payment.GetPaymentId(),
		Status:    payment.GetStatus().String(),
	}, nil
}

func (s *PaymentService) GenerateReceipt(ctx context.Context, req *paymentservicev1.GenerateReceiptRequest) (*paymentservicev1.GenerateReceiptResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}
	payment, err := s.repo.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, mapRepoError(err)
	}
	if !isSettledStatus(payment.GetStatus()) {
		return nil, fmt.Errorf("%w: receipt requires a settled payment", ErrInvalidTransition)
	}

	now := time.Now().UTC()
	s.ensureReceiptFields(payment, now)
	payment.UpdatedAt = timestamppb.New(now)
	payment.ReceiptUrl = s.receiptURL(payment.GetPaymentId())
	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}
	return &paymentservicev1.GenerateReceiptResponse{
		PaymentId:     payment.GetPaymentId(),
		ReceiptNumber: payment.GetReceiptNumber(),
		ReceiptFileId: payment.GetReceiptFileId(),
	}, nil
}

func (s *PaymentService) GetPaymentReceipt(ctx context.Context, req *paymentservicev1.GetPaymentReceiptRequest) (*paymentservicev1.GetPaymentReceiptResponse, error) {
	if req == nil || strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}
	payment, err := s.repo.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, mapRepoError(err)
	}
	if payment.GetReceiptNumber() == "" && isSettledStatus(payment.GetStatus()) {
		now := time.Now().UTC()
		s.ensureReceiptFields(payment, now)
		payment.UpdatedAt = timestamppb.New(now)
		payment.ReceiptUrl = s.receiptURL(payment.GetPaymentId())
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
	}
	if payment.GetReceiptNumber() == "" {
		return nil, fmt.Errorf("%w: receipt has not been generated", ErrInvalidTransition)
	}
	return &paymentservicev1.GetPaymentReceiptResponse{
		PaymentId:     payment.GetPaymentId(),
		ReceiptNumber: payment.GetReceiptNumber(),
		ReceiptFileId: payment.GetReceiptFileId(),
		ReceiptUrl:    payment.GetReceiptUrl(),
		GeneratedAt:   firstNonNilTimestamp(payment.GetUpdatedAt(), payment.GetCompletedAt(), payment.GetCreatedAt()),
	}, nil
}

func (s *PaymentService) buildGatewayURLs(paymentID, callbackURL string) (string, string, string, string) {
	base := strings.TrimRight(s.config.PublicBaseURL, "/")
	if base == "" {
		base = callbackBaseURL(callbackURL)
	}
	if base == "" {
		base = "http://localhost"
	}
	return base + "/v1/payments/sslcommerz/success?payment_id=" + paymentID,
		base + "/v1/payments/sslcommerz/fail?payment_id=" + paymentID,
		base + "/v1/payments/sslcommerz/cancel?payment_id=" + paymentID,
		base + "/v1/payments/webhook/sslcommerz?payment_id=" + paymentID
}

func callbackBaseURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}

func callbackFields(values url.Values) map[string]string {
	fields := map[string]string{}
	for key, list := range values {
		if len(list) == 0 {
			continue
		}
		fields[key] = strings.TrimSpace(list[0])
	}
	return fields
}

func (s *PaymentService) lookupPaymentForCallback(ctx context.Context, provider string, values url.Values) (*paymententityv1.Payment, error) {
	if paymentID := firstNonEmpty(values.Get("value_a"), values.Get("payment_id")); paymentID != "" {
		return s.repo.GetPayment(ctx, paymentID)
	}
	for _, reference := range []string{values.Get("val_id"), values.Get("tran_id"), values.Get("sessionkey")} {
		if strings.TrimSpace(reference) == "" {
			continue
		}
		payment, err := s.repo.GetPaymentByProviderReference(ctx, provider, strings.TrimSpace(reference))
		if err == nil {
			return payment, nil
		}
		if err != domain.ErrNotFound {
			return nil, err
		}
	}
	if orderID := strings.TrimSpace(values.Get("value_b")); orderID != "" {
		return s.repo.GetPaymentByOrderID(ctx, orderID)
	}
	return nil, fmt.Errorf("%w: unable to resolve payment from provider callback", ErrNotFound)
}

func (s *PaymentService) ensureReceiptFields(payment *paymententityv1.Payment, now time.Time) {
	if payment == nil {
		return
	}
	if strings.TrimSpace(payment.GetReceiptNumber()) == "" {
		payment.ReceiptNumber = fmt.Sprintf("RCP-%s-%s", now.Format("20060102"), strings.ToUpper(strings.ReplaceAll(payment.GetPaymentId(), "-", ""))[:8])
	}
	if strings.TrimSpace(payment.GetReceiptUrl()) == "" {
		payment.ReceiptUrl = s.receiptURL(payment.GetPaymentId())
	}
}

func (s *PaymentService) receiptURL(paymentID string) string {
	base := strings.TrimRight(s.config.PublicBaseURL, "/")
	if base == "" {
		return "/v1/payments/" + paymentID + "/receipt"
	}
	return base + "/v1/payments/" + paymentID + "/receipt"
}

func firstNonNilTimestamp(values ...*timestamppb.Timestamp) *timestamppb.Timestamp {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func normalizeMethod(value string) (paymententityv1.PaymentMethod, string, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "BANK_TRANSFER":
		return paymententityv1.PaymentMethod_PAYMENT_METHOD_BANK_TRANSFER, "BANK_TRANSFER", nil
	case "CARD", "SSLCOMMERZ":
		return paymententityv1.PaymentMethod_PAYMENT_METHOD_CARD, "SSLCOMMERZ", nil
	default:
		return paymententityv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, "", fmt.Errorf("%w: unsupported payment_method", ErrNotImplemented)
	}
}

func resolveUserID(ctx context.Context, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	// Only read from the authoritative gateway-set header — x-user-id.
	// x-customer-id and x-subject are NOT set by the gateway and must not be used.
	if values := md.Get("x-user-id"); len(values) > 0 && strings.TrimSpace(values[0]) != "" {
		return strings.TrimSpace(values[0])
	}
	return ""
}

// paymentRequestContext is a lightweight caller-identity snapshot read from gRPC metadata.
// Used by publish helpers to stamp events with full context without re-exposing the middleware package.
type paymentRequestContext struct {
	userID    string
	tenantID  string
	orgID     string
	portal    string
	sessionID string
	tokenID   string
}

// extractPaymentContext reads the gateway-set metadata headers from ctx.
// Returns zero-value struct if metadata is absent (e.g. webhook callbacks, unit tests).
func extractPaymentContext(ctx context.Context) paymentRequestContext {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return paymentRequestContext{}
	}
	first := func(key string) string {
		for _, v := range md.Get(key) {
			if v = strings.TrimSpace(v); v != "" {
				return v
			}
		}
		return ""
	}
	rawPortal := first("x-portal")
	portal := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(rawPortal), "PORTAL_"))
	return paymentRequestContext{
		userID:    first("x-user-id"),
		tenantID:  first("x-tenant-id"),
		orgID:     first("x-business-id"),
		portal:    portal,
		sessionID: first("x-session-id"),
		tokenID:   first("x-token-id"),
	}
}

func referenceInfo(req *paymentservicev1.InitiatePaymentRequest) (string, string) {
	if req.GetPolicyId() != "" {
		return "policy", req.GetPolicyId()
	}
	if req.GetMetadata() != nil {
		referenceType := strings.TrimSpace(req.GetMetadata()["reference_type"])
		referenceID := strings.TrimSpace(req.GetMetadata()["reference_id"])
		if referenceType != "" || referenceID != "" {
			return referenceType, referenceID
		}
	}
	return "", ""
}

func parsePageToken(value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}
	var offset int
	_, err := fmt.Sscanf(value, "%d", &offset)
	return offset, err
}

func normalizePageSize(value int32) int32 {
	if value <= 0 {
		return 20
	}
	if value > 100 {
		return 100
	}
	return value
}

func tsPtr(value *timestamppb.Timestamp) *time.Time {
	if value == nil {
		return nil
	}
	t := value.AsTime().UTC()
	return &t
}

func isTerminalFailure(status paymententityv1.PaymentStatus) bool {
	switch status {
	case paymententityv1.PaymentStatus_PAYMENT_STATUS_FAILED,
		paymententityv1.PaymentStatus_PAYMENT_STATUS_CANCELLED,
		paymententityv1.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return true
	default:
		return false
	}
}

func isSettledStatus(status paymententityv1.PaymentStatus) bool {
	switch status {
	case paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS,
		paymententityv1.PaymentStatus_PAYMENT_STATUS_VERIFIED,
		paymententityv1.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return true
	default:
		return false
	}
}

func mapRepoError(err error) error {
	if err == nil {
		return nil
	}
	if err == domain.ErrNotFound {
		return ErrNotFound
	}
	return err
}

func cloneMoney(value *commonv1.Money, currency string) *commonv1.Money {
	if value == nil {
		return nil
	}
	// Construct a fresh proto message — never copy a proto struct by value
	// because proto messages embed sync.Mutex via protoimpl.MessageState.
	cur := value.GetCurrency()
	if cur == "" {
		cur = currency
	}
	amt := value.GetAmount()
	dec := value.GetDecimalAmount()
	if dec == 0 && amt != 0 {
		dec = float64(amt) / 100
	}
	return &commonv1.Money{
		Amount:        amt,
		Currency:      cur,
		DecimalAmount: dec,
	}
}

func gatewayValue(payload string, key string) string {
	values := paymentrepo.UnmarshalGatewayResponse(payload)
	if values == nil {
		return ""
	}
	return values[key]
}

func mergeGatewayResponse(existing string, updates map[string]string) string {
	payload := paymentrepo.UnmarshalGatewayResponse(existing)
	if payload == nil {
		payload = map[string]string{}
	}
	for key, value := range updates {
		if strings.TrimSpace(value) != "" {
			payload[key] = value
		}
	}
	return paymentrepo.MarshalGatewayResponse(payload)
}

func newError(code string, message string) *commonv1.Error {
	return &commonv1.Error{Code: code, Message: message}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func merchantTransactionID(paymentID string) string {
	return "tran-" + strings.ReplaceAll(paymentID, "-", "")
}

func normalizeOptionalUUID(value string) string {
	if _, err := uuid.Parse(strings.TrimSpace(value)); err == nil {
		return strings.TrimSpace(value)
	}
	return ""
}

func providerStatusValid(resp *domain.GatewayValidationResponse) bool {
	if resp == nil {
		return false
	}
	switch strings.ToUpper(strings.TrimSpace(resp.Status)) {
	case "VALID", "VALIDATED", "SUCCESS":
		return true
	default:
		return false
	}
}

func validationMatches(payment *paymententityv1.Payment, resp *domain.GatewayValidationResponse) bool {
	if payment == nil || resp == nil {
		return false
	}
	if txn := strings.TrimSpace(resp.TransactionID); txn != "" && payment.GetTransactionId() != "" && payment.GetTransactionId() != txn {
		return false
	}
	if resp.Amount != nil && payment.GetAmount() != nil {
		if resp.Amount.GetAmount() != 0 && payment.GetAmount().GetAmount() != resp.Amount.GetAmount() {
			return false
		}
		if resp.Amount.GetCurrency() != "" && payment.GetCurrency() != "" && !strings.EqualFold(payment.GetCurrency(), resp.Amount.GetCurrency()) {
			return false
		}
	}
	return true
}

func (s *PaymentService) publishInitiated(ctx context.Context, payment *paymententityv1.Payment, refType, refID string) {
	if s.publisher == nil || payment == nil {
		return
	}
	rctx := extractPaymentContext(ctx)
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicPaymentInitiated, payment.GetPaymentId(), &paymenteventsv1.PaymentInitiatedEvent{
		EventId:       uuid.NewString(),
		PaymentId:     payment.GetPaymentId(),
		PayerId:       payment.GetPayerId(),
		Amount:        payment.GetAmount(),
		PaymentMethod: payment.GetMethod().String(),
		ReferenceType: firstNonEmpty(refType, "order"),
		ReferenceId:   firstNonEmpty(refID, payment.GetOrderId()),
		Timestamp:     now,
		CorrelationId: payment.GetOrderId(),
		// Full typed context fields
		OrderId:        payment.GetOrderId(),
		InvoiceId:      payment.GetInvoiceId(),
		TenantId:       firstNonEmpty(payment.GetTenantId(), rctx.tenantID),
		OrganisationId: firstNonEmpty(payment.GetOrganisationId(), rctx.orgID),
		Provider:       payment.GetProvider(),
		OccurredAt:     now,
	})
}

func (s *PaymentService) publishCompleted(ctx context.Context, payment *paymententityv1.Payment, correlationID string) {
	if s.publisher == nil || payment == nil {
		return
	}
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicPaymentCompleted, payment.GetPaymentId(), &paymenteventsv1.PaymentCompletedEvent{
		EventId:       uuid.NewString(),
		PaymentId:     payment.GetPaymentId(),
		TransactionId: payment.GetTransactionId(),
		Amount:        payment.GetAmount(),
		PaymentMethod: payment.GetMethod().String(),
		Timestamp:     now,
		CorrelationId: firstNonEmpty(correlationID, payment.GetOrderId()),
		// Extended typed fields
		OrderId:        payment.GetOrderId(),
		InvoiceId:      payment.GetInvoiceId(),
		TenantId:       payment.GetTenantId(),
		OrganisationId: payment.GetOrganisationId(),
		Provider:       payment.GetProvider(),
		ValId:          payment.GetValId(),
		ReceiptNumber:  payment.GetReceiptNumber(),
		OccurredAt:     now,
	})
}

func (s *PaymentService) publishFailed(ctx context.Context, payment *paymententityv1.Payment, code, message string) {
	if s.publisher == nil || payment == nil {
		return
	}
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicPaymentFailed, payment.GetPaymentId(), &paymenteventsv1.PaymentFailedEvent{
		EventId:       uuid.NewString(),
		PaymentId:     payment.GetPaymentId(),
		PayerId:       payment.GetPayerId(),
		Amount:        payment.GetAmount(),
		PaymentMethod: payment.GetMethod().String(),
		ErrorCode:     code,
		ErrorMessage:  message,
		Timestamp:     now,
		CorrelationId: payment.GetOrderId(),
		// Extended typed fields
		OrderId:        payment.GetOrderId(),
		InvoiceId:      payment.GetInvoiceId(),
		TenantId:       payment.GetTenantId(),
		OrganisationId: payment.GetOrganisationId(),
		Provider:       payment.GetProvider(),
		OccurredAt:     now,
	})
}

func (s *PaymentService) publishRefundProcessed(ctx context.Context, refund *paymententityv1.PaymentRefund, payment *paymententityv1.Payment) {
	if s.publisher == nil || refund == nil || payment == nil {
		return
	}
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicRefundProcessed, refund.GetRefundId(), &paymenteventsv1.RefundProcessedEvent{
		EventId:           uuid.NewString(),
		RefundId:          refund.GetRefundId(),
		OriginalPaymentId: refund.GetPaymentId(),
		RecipientId:       payment.GetPayerId(),
		Amount:            refund.GetRefundAmount(),
		Reason:            refund.GetReason(),
		Timestamp:         now,
		CorrelationId:     payment.GetOrderId(),
		// Extended typed fields
		OrderId:        payment.GetOrderId(),
		InvoiceId:      payment.GetInvoiceId(),
		TenantId:       payment.GetTenantId(),
		OrganisationId: payment.GetOrganisationId(),
		OccurredAt:     now,
	})
}

// publishManualProofSubmitted publishes insuretech.payment.v1.payment.manual_review_requested.
// Called by SubmitManualPaymentProof after the proof file_id has been persisted.
func (s *PaymentService) publishManualProofSubmitted(ctx context.Context, payment *paymententityv1.Payment, submittedBy string) {
	if s.publisher == nil || payment == nil {
		return
	}
	rctx := extractPaymentContext(ctx)
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicManualReviewRequested, payment.GetPaymentId(), &paymenteventsv1.ManualPaymentProofSubmittedEvent{
		EventId:        uuid.NewString(),
		PaymentId:      payment.GetPaymentId(),
		OrderId:        payment.GetOrderId(),
		SubmittedBy:    firstNonEmpty(submittedBy, rctx.userID),
		TenantId:       firstNonEmpty(payment.GetTenantId(), rctx.tenantID),
		OrganisationId: firstNonEmpty(payment.GetOrganisationId(), rctx.orgID),
		OccurredAt:     now,
	})
}

// publishVerified publishes insuretech.payment.v1.payment.verified.
// Called by ReviewManualPayment when approved=true.
func (s *PaymentService) publishVerified(ctx context.Context, payment *paymententityv1.Payment, verifiedBy string) {
	if s.publisher == nil || payment == nil {
		return
	}
	rctx := extractPaymentContext(ctx)
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicPaymentVerified, payment.GetPaymentId(), &paymenteventsv1.PaymentVerifiedEvent{
		EventId:        uuid.NewString(),
		PaymentId:      payment.GetPaymentId(),
		OrderId:        payment.GetOrderId(),
		VerifiedBy:     firstNonEmpty(verifiedBy, rctx.userID),
		TenantId:       firstNonEmpty(payment.GetTenantId(), rctx.tenantID),
		OrganisationId: firstNonEmpty(payment.GetOrganisationId(), rctx.orgID),
		CorrelationId:  payment.GetIdempotencyKey(),
		CausationId:    payment.GetPaymentId(),
		Amount:         payment.GetAmount(),
		VerifiedAt:     now,
		OccurredAt:     now,
	})
}

// publishManualReviewed publishes insuretech.payment.v1.payment.manual_review_completed.
func (s *PaymentService) publishManualReviewed(ctx context.Context, payment *paymententityv1.Payment, reviewedBy, notes, rejectionReason string, approved bool) {
	if s.publisher == nil || payment == nil {
		return
	}
	rctx := extractPaymentContext(ctx)
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicManualPaymentReviewed, payment.GetPaymentId(), &paymenteventsv1.ManualPaymentReviewedEvent{
		EventId:         uuid.NewString(),
		PaymentId:       payment.GetPaymentId(),
		OrderId:         payment.GetOrderId(),
		Approved:        approved,
		ReviewedBy:      firstNonEmpty(reviewedBy, rctx.userID),
		ReviewNotes:     notes,
		RejectionReason: rejectionReason,
		TenantId:        firstNonEmpty(payment.GetTenantId(), rctx.tenantID),
		OrganisationId:  firstNonEmpty(payment.GetOrganisationId(), rctx.orgID),
		OccurredAt:      now,
	})
}

// publishReceiptGenerated publishes insuretech.payment.v1.payment.receipt_generated.
func (s *PaymentService) publishReceiptGenerated(ctx context.Context, payment *paymententityv1.Payment) {
	if s.publisher == nil || payment == nil {
		return
	}
	rctx := extractPaymentContext(ctx)
	now := timestamppb.Now()
	s.publisher.Publish(ctx, events.TopicReceiptGenerated, payment.GetPaymentId(), &paymenteventsv1.ReceiptGeneratedEvent{
		EventId:        uuid.NewString(),
		PaymentId:      payment.GetPaymentId(),
		OrderId:        payment.GetOrderId(),
		ReceiptNumber:  payment.GetReceiptNumber(),
		TenantId:       firstNonEmpty(payment.GetTenantId(), rctx.tenantID),
		OrganisationId: firstNonEmpty(payment.GetOrganisationId(), rctx.orgID),
		OccurredAt:     now,
	})
}

// ─── Phase 2 RPC implementations ─────────────────────────────────────────────

