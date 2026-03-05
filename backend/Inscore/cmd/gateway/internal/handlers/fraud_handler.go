package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	fraudservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/fraud/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FraudHandler proxies fraud APIs to the fraud gRPC service.
type FraudHandler struct {
	client FraudClient
}

// FraudClient abstracts fraud gRPC client methods used by gateway.
type FraudClient interface {
	CheckFraud(ctx context.Context, in *fraudservicev1.CheckFraudRequest, opts ...grpc.CallOption) (*fraudservicev1.CheckFraudResponse, error)
	GetFraudAlert(ctx context.Context, in *fraudservicev1.GetFraudAlertRequest, opts ...grpc.CallOption) (*fraudservicev1.GetFraudAlertResponse, error)
	ListFraudAlerts(ctx context.Context, in *fraudservicev1.ListFraudAlertsRequest, opts ...grpc.CallOption) (*fraudservicev1.ListFraudAlertsResponse, error)
	CreateFraudCase(ctx context.Context, in *fraudservicev1.CreateFraudCaseRequest, opts ...grpc.CallOption) (*fraudservicev1.CreateFraudCaseResponse, error)
	GetFraudCase(ctx context.Context, in *fraudservicev1.GetFraudCaseRequest, opts ...grpc.CallOption) (*fraudservicev1.GetFraudCaseResponse, error)
	UpdateFraudCase(ctx context.Context, in *fraudservicev1.UpdateFraudCaseRequest, opts ...grpc.CallOption) (*fraudservicev1.UpdateFraudCaseResponse, error)
	ListFraudRules(ctx context.Context, in *fraudservicev1.ListFraudRulesRequest, opts ...grpc.CallOption) (*fraudservicev1.ListFraudRulesResponse, error)
	CreateFraudRule(ctx context.Context, in *fraudservicev1.CreateFraudRuleRequest, opts ...grpc.CallOption) (*fraudservicev1.CreateFraudRuleResponse, error)
	UpdateFraudRule(ctx context.Context, in *fraudservicev1.UpdateFraudRuleRequest, opts ...grpc.CallOption) (*fraudservicev1.UpdateFraudRuleResponse, error)
	ActivateFraudRule(ctx context.Context, in *fraudservicev1.ActivateFraudRuleRequest, opts ...grpc.CallOption) (*fraudservicev1.ActivateFraudRuleResponse, error)
	DeactivateFraudRule(ctx context.Context, in *fraudservicev1.DeactivateFraudRuleRequest, opts ...grpc.CallOption) (*fraudservicev1.DeactivateFraudRuleResponse, error)
}

func NewFraudHandler(conn *grpc.ClientConn) *FraudHandler {
	return &FraudHandler{client: fraudservicev1.NewFraudServiceClient(conn)}
}

func (h *FraudHandler) Check(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req fraudservicev1.CheckFraudRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CheckFraud(ctx, &req)
	})
}

func (h *FraudHandler) GetAlert(w http.ResponseWriter, r *http.Request) {
	alertID := r.PathValue("fraud_alert_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetFraudAlert(ctx, &fraudservicev1.GetFraudAlertRequest{FraudAlertId: alertID})
	})
}

func (h *FraudHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &fraudservicev1.ListFraudAlertsRequest{
			Status:    r.URL.Query().Get("status"),
			RiskLevel: r.URL.Query().Get("risk_level"),
			StartDate: r.URL.Query().Get("start_date"),
			EndDate:   r.URL.Query().Get("end_date"),
		}
		if q := r.URL.Query().Get("page"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.Page = int32(n)
			}
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListFraudAlerts(ctx, req)
	})
}

func (h *FraudHandler) CreateCase(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req fraudservicev1.CreateFraudCaseRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateFraudCase(ctx, &req)
	})
}

func (h *FraudHandler) GetCase(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("fraud_case_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetFraudCase(ctx, &fraudservicev1.GetFraudCaseRequest{FraudCaseId: caseID})
	})
}

func (h *FraudHandler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("fraud_case_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req fraudservicev1.UpdateFraudCaseRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.FraudCaseId = caseID
		return h.client.UpdateFraudCase(ctx, &req)
	})
}

func (h *FraudHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &fraudservicev1.ListFraudRulesRequest{
			Category:  r.URL.Query().Get("category"),
			PageToken: r.URL.Query().Get("page_token"),
		}
		if q := r.URL.Query().Get("active_only"); q == "1" || q == "true" {
			req.ActiveOnly = true
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListFraudRules(ctx, req)
	})
}

func (h *FraudHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req fraudservicev1.CreateFraudRuleRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateFraudRule(ctx, &req)
	})
}

func (h *FraudHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	ruleID := r.PathValue("rule_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req fraudservicev1.UpdateFraudRuleRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.RuleId = ruleID
		return h.client.UpdateFraudRule(ctx, &req)
	})
}

func (h *FraudHandler) ActivateRule(w http.ResponseWriter, r *http.Request) {
	ruleID := r.PathValue("rule_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req fraudservicev1.ActivateFraudRuleRequest
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, &req); err != nil {
				return nil, err
			}
		}
		req.RuleId = ruleID
		return h.client.ActivateFraudRule(ctx, &req)
	})
}

func (h *FraudHandler) DeactivateRule(w http.ResponseWriter, r *http.Request) {
	ruleID := r.PathValue("rule_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req fraudservicev1.DeactivateFraudRuleRequest
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, &req); err != nil {
				return nil, err
			}
		}
		req.RuleId = ruleID
		return h.client.DeactivateFraudRule(ctx, &req)
	})
}

func parseHTTPTime(raw string) (*timestamppb.Timestamp, error) {
	if raw == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return timestamppb.New(t), nil
		}
	}
	return nil, errors.New("invalid timestamp format")
}
