package repository

import (
	"fmt"
	"strings"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var mediaProtoJSONMarshaler = protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: false}
var mediaProtoJSONUnmarshaler = protojson.UnmarshalOptions{DiscardUnknown: true}

func prepareAuditInfoForCreate(info *commonv1.AuditInfo, defaultActor string) (*commonv1.AuditInfo, string, error) {
	var cloned *commonv1.AuditInfo
	if info != nil {
		msg, ok := proto.Clone(info).(*commonv1.AuditInfo)
		if !ok {
			return nil, "", fmt.Errorf("clone audit info: unexpected type %T", info)
		}
		cloned = msg
	} else {
		cloned = &commonv1.AuditInfo{}
	}

	now := timestamppb.Now()
	actor := strings.TrimSpace(defaultActor)

	if cloned.CreatedAt == nil {
		cloned.CreatedAt = now
	}
	if cloned.UpdatedAt == nil {
		cloned.UpdatedAt = now
	}
	if cloned.CreatedBy == "" {
		cloned.CreatedBy = actor
	}
	if cloned.UpdatedBy == "" {
		cloned.UpdatedBy = cloned.CreatedBy
	}

	raw, err := mediaProtoJSONMarshaler.Marshal(cloned)
	if err != nil {
		return nil, "", fmt.Errorf("marshal audit info: %w", err)
	}

	return cloned, string(raw), nil
}

func parseAuditInfoJSON(raw string) (*commonv1.AuditInfo, error) {
	info := &commonv1.AuditInfo{}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "null" {
		return info, nil
	}
	if err := mediaProtoJSONUnmarshaler.Unmarshal([]byte(trimmed), info); err != nil {
		return nil, fmt.Errorf("unmarshal audit info: %w", err)
	}
	return info, nil
}

func auditInfoSelectExpr(alias string) string {
	if alias != "" {
		alias += "."
	}
	return fmt.Sprintf("COALESCE(%saudit_info::text, '{}')", alias)
}

func auditCreatedAtExpr(alias string) string {
	if alias != "" {
		alias += "."
	}
	return fmt.Sprintf("COALESCE((%saudit_info->>'created_at')::timestamptz, %sstarted_at, %scompleted_at, NOW())", alias, alias, alias)
}

func auditCreatedAtExprMedia(alias string) string {
	if alias != "" {
		alias += "."
	}
	return fmt.Sprintf("COALESCE((%saudit_info->>'created_at')::timestamptz, NOW())", alias)
}

func auditInfoUpdatedExpr(alias string) string {
	if alias != "" {
		alias += "."
	}
	return fmt.Sprintf("%saudit_info = COALESCE(%saudit_info, '{}'::jsonb) || jsonb_build_object('updated_at', to_jsonb(NOW()))", alias, alias)
}
