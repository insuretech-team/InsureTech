using System;
using System.Collections.Generic;
using InsuranceEngine.Fraud.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Fraud.Application.Features.Commands.CheckFraud;

    Guid ClaimId, 
    Guid PolicyId, 
    Guid CustomerId,
    long ClaimedAmount, 
    long SumInsuredAmount,
    string ClaimType,
    string? PlaceOfIncident,
    DateTime IncidentDate, 
    DateTime PolicyIssuedAt) : IRequest<Result<FraudCheckResponse>>;

public record FraudCheckResponse(
    Guid CheckId, 
    FraudRiskLevel RiskLevel, 
    FraudCheckStatus Status, 
    double RiskScore, 
    List<string> Findings);
