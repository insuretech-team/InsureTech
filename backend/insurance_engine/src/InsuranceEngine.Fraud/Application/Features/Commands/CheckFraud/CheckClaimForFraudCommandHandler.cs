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

        // FD-001: Rapid Policy-Claim (within 7 days)
        var daysSinceIssuance = (request.IncidentDate - request.PolicyIssuedAt).TotalDays;
        if (daysSinceIssuance <= 7)
        {
            findings.Add("FD-001: Claim submitted within 7 days of policy issuance.");
            riskScore += 40;
        }
        else if (daysSinceIssuance <= 30)
        {
            findings.Add("FD-001: Claim submitted within 30 days of policy issuance.");
            riskScore += 20;
        }

        // FD-002: High relative amount (Placeholder logic)
        // In a real system, we'd check Sum Insured.
        if (request.ClaimedAmount > 10_000_000) // > 100k BDT
        {
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
