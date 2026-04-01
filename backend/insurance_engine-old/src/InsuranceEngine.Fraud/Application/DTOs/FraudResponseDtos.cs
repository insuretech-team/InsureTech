using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Fraud.Application.DTOs;

public record CheckFraudRequest(
    [property: JsonPropertyName("entity_type")] string EntityType,
    [property: JsonPropertyName("entity_id")] string EntityId,
    [property: JsonPropertyName("data")] object Data
);

public record CheckFraudResponse(
    [property: JsonPropertyName("is_fraud_detected")] bool IsFraudDetected,
    [property: JsonPropertyName("fraud_score")] int FraudScore,
    [property: JsonPropertyName("risk_level")] string RiskLevel,
    [property: JsonPropertyName("triggered_rules")] List<string> TriggeredRules,
    [property: JsonPropertyName("fraud_alert_id")] string FraudAlertId,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record FraudAlertsListingResponse(
    [property: JsonPropertyName("fraud_alerts")] List<FraudAlertDto> FraudAlerts,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record FraudAlertDto(
    [property: JsonPropertyName("fraud_alert_id")] string FraudAlertId,
    [property: JsonPropertyName("entity_type")] string EntityType,
    [property: JsonPropertyName("entity_id")] string EntityId,
    [property: JsonPropertyName("fraud_score")] int FraudScore,
    [property: JsonPropertyName("risk_level")] string RiskLevel,
    [property: JsonPropertyName("triggered_rules")] List<string> TriggeredRules,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("created_at")] DateTime CreatedAt
);

public record FraudAlertRetrievalResponse(
    [property: JsonPropertyName("fraud_alert")] FraudAlertDto FraudAlert,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record FraudCasesListingResponse(
    [property: JsonPropertyName("fraud_cases")] List<FraudCaseDto> FraudCases,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record FraudCaseDto(
    [property: JsonPropertyName("fraud_case_id")] string FraudCaseId,
    [property: JsonPropertyName("fraud_alert_id")] string FraudAlertId,
    [property: JsonPropertyName("investigator_id")] string? InvestigatorId,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("priority")] string Priority,
    [property: JsonPropertyName("notes")] string? Notes,
    [property: JsonPropertyName("created_at")] DateTime CreatedAt,
    [property: JsonPropertyName("updated_at")] DateTime UpdatedAt
);

public record FraudCaseRetrievalResponse(
    [property: JsonPropertyName("fraud_case")] FraudCaseDto FraudCase,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record FraudRulesListingResponse(
    [property: JsonPropertyName("fraud_rules")] List<FraudRuleDto> FraudRules,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record FraudRuleDto(
    [property: JsonPropertyName("rule_id")] string RuleId,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("description")] string Description,
    [property: JsonPropertyName("entity_type")] string EntityType,
    [property: JsonPropertyName("criteria")] string Criteria,
    [property: JsonPropertyName("is_active")] bool IsActive,
    [property: JsonPropertyName("created_at")] DateTime CreatedAt
);

public record FraudRuleRetrievalResponse(
    [property: JsonPropertyName("fraud_rule")] FraudRuleDto FraudRule,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record FraudRuleOperationResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
