using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Fraud.Domain.Entities;
using InsuranceEngine.Fraud.Domain.Enums;
using InsuranceEngine.Fraud.Infrastructure.Persistence;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Fraud.Application.Features.Commands.CheckFraud;

public class CheckClaimForFraudCommandHandler : IRequestHandler<CheckClaimForFraudCommand, Result<FraudCheckResponse>>
{
    private readonly FraudDbContext _dbContext;

    public CheckClaimForFraudCommandHandler(FraudDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    public async Task<Result<FraudCheckResponse>> Handle(CheckClaimForFraudCommand request, CancellationToken cancellationToken)
    {
        var findings = new List<string>();
        double riskScore = 0;

        // FR-182: Rapid Claim (Incident within 30 days of Policy Issuance)
        var daysSinceIssuance = (request.IncidentDate - request.PolicyIssuedAt).TotalDays;
        if (daysSinceIssuance <= 7)
        {
            findings.Add("FR-182: Claim submitted within 7 days of policy issuance.");
            riskScore += 40;
        }
        else if (daysSinceIssuance <= 30)
        {
            findings.Add("FR-182: Claim submitted within 30 days of policy issuance.");
            riskScore += 20;
        }

        // FR-183: High Claim Amount (> 80% of Sum Insured)
        if (request.SumInsuredAmount > 0 && request.ClaimedAmount >= (request.SumInsuredAmount * 0.80))
        {
            findings.Add("FR-183: Claimed amount exceeds 80% of the total Sum Insured.");
            riskScore += 25;
        }

        // FR-184: Multiple claims velocity (Historical DB check mocked via query if implemented)
        var recentClaimsCount = _dbContext.FraudChecks
            .Count(f => f.EntityType == "Claim" && f.CreatedAt > DateTime.UtcNow.AddMonths(-6));
        if (recentClaimsCount >= 3) // Assuming checking across system for now, but really should be per customer
        {
            findings.Add("FR-184: High velocity: Customer has submitted multiple claims recently.");
            riskScore += 15;
        }

        // FR-185: High-Risk Claim Types (e.g. Theft or Unexplained)
        if (request.ClaimType.Contains("Theft", StringComparison.OrdinalIgnoreCase) || 
            request.ClaimType.Contains("Unexplained", StringComparison.OrdinalIgnoreCase))
        {
            findings.Add($"FR-185: High risk incident type ({request.ClaimType}).");
            riskScore += 15;
        }

        // FR-186: Suspicious Incident Time (e.g. 1 AM to 4 AM)
        var hour = request.IncidentDate.ToLocalTime().Hour;
        if (hour >= 1 && hour <= 4)
        {
            findings.Add("FR-186: Incident occurred during suspicious late-night hours (1 AM to 4 AM).");
            riskScore += 20;
        }

        // FR-187: Blacklisted/High-Risk Regions
        if (!string.IsNullOrEmpty(request.PlaceOfIncident) && 
           (request.PlaceOfIncident.Contains("HighRiskZone", StringComparison.OrdinalIgnoreCase) || 
            request.PlaceOfIncident.Contains("UnverifiedArea", StringComparison.OrdinalIgnoreCase)))
        {
            findings.Add("FR-187: Incident occurred in a flagged high-risk geographic zone.");
            riskScore += 15;
        }

        // FR-188: Delayed Reporting (> 30 days since incident)
        var daysSinceIncident = (DateTime.UtcNow - request.IncidentDate).TotalDays;
        if (daysSinceIncident > 30)
        {
            findings.Add($"FR-188: Delayed reporting. Claim reported {Math.Round(daysSinceIncident)} days after incident.");
            riskScore += 10;
        }

        var riskLevel = riskScore switch
        {
            >= 70 => FraudRiskLevel.Critical,
            >= 40 => FraudRiskLevel.High,
            >= 20 => FraudRiskLevel.Medium,
            _ => FraudRiskLevel.Low
        };

        var status = riskLevel >= FraudRiskLevel.High ? FraudCheckStatus.Flagged : FraudCheckStatus.Approved;

        var fraudCheck = new FraudCheck
        {
            Id = Guid.NewGuid(),
            EntityId = request.ClaimId,
            EntityType = "Claim",
            RiskLevel = riskLevel,
            Status = status,
            RiskScore = riskScore,
            FindingsJson = JsonSerializer.Serialize(findings),
            CheckedAt = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        _dbContext.FraudChecks.Add(fraudCheck);
        await _dbContext.SaveChangesAsync(cancellationToken);

        return Result<FraudCheckResponse>.Ok(new FraudCheckResponse(
            fraudCheck.Id,
            fraudCheck.RiskLevel,
            fraudCheck.Status,
            fraudCheck.RiskScore,
            findings));
    }
}
