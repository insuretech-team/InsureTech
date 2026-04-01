using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Fraud.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.FraudDetection.Application.Commands;

public sealed class CheckFraudCommandHandler : IRequestHandler<CheckFraudCommand, CheckFraudResponse>
{
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly ILogger<CheckFraudCommandHandler> _logger;

    public CheckFraudCommandHandler(
        IRepository<PolicyEntity> policyRepository,
        IRepository<ClaimEntity> claimRepository,
        ILogger<CheckFraudCommandHandler> logger)
    {
        _policyRepository = policyRepository;
        _claimRepository = claimRepository;
        _logger = logger;
    }

    public async Task<CheckFraudResponse> Handle(CheckFraudCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var response = new CheckFraudResponse
            {
                IsFraudDetected = false,
                FraudScore = 0,
                RiskLevel = "LOW"
            };

            // FR-094/TM-001/SEC-022: AML & Fraud Rule Checks
            
            if (request.EntityType == "POLICY")
            {
                await PerformPolicyFraudCheck(request, response, cancellationToken);
            }
            else if (request.EntityType == "CLAIM")
            {
                await PerformClaimFraudCheck(request, response, cancellationToken);
            }

            // Determine final risk level based on score
            if (response.FraudScore >= 70) response.RiskLevel = "CRITICAL";
            else if (response.FraudScore >= 40) response.RiskLevel = "MEDIUM";
            else response.RiskLevel = "LOW";

            if (response.FraudScore > 0)
            {
                _logger.LogWarning("Fraud Check for {EntityType} {EntityId}: Score {Score}, Risk {Risk}", 
                    request.EntityType, request.EntityId, response.FraudScore, response.RiskLevel);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Fraud check failed for {EntityId}", request.EntityId);
            return new CheckFraudResponse
            {
                Error = new Error { Code = "FRAUD_CHECK_FAILED", Message = ex.Message }
            };
        }
    }

    private async Task PerformPolicyFraudCheck(CheckFraudCommand request, CheckFraudResponse response, CancellationToken cancellationToken)
    {
        // SEC-022: Premium > BDT 5L without income proof (simulated check)
        if (request.Data.Fields.TryGetValue("premium_amount", out var premium) && premium.NumberValue > 500000)
        {
            if (!request.Data.Fields.TryGetValue("has_income_proof", out var incomeProof) || !incomeProof.BoolValue)
            {
                response.IsFraudDetected = true;
                response.FraudScore += 30;
                response.TriggeredRules.Add("SEC-022: High premium without income proof");
            }
        }

        // TM-001: Multiple policies (>3) in 7 days for same customer (simulated check)
        if (request.Data.Fields.TryGetValue("customer_id", out var customerId))
        {
            var userId = Guid.Parse(customerId.StringValue);
            var sevenDaysAgo = DateTime.UtcNow.AddDays(-7);
            var recentPolicies = await _policyRepository.FindAsync(p => p.CustomerId == userId && p.CreatedAt >= sevenDaysAgo, cancellationToken);
            
            if (recentPolicies.Count() >= 3)
            {
                response.IsFraudDetected = true;
                response.FraudScore += 40;
                response.TriggeredRules.Add("TM-001: Excessive policies (>3) in 7 days");
            }
        }
    }

    private async Task PerformClaimFraudCheck(CheckFraudCommand request, CheckFraudResponse response, CancellationToken cancellationToken)
    {
        // FR-094: Frequent claims (>3 in 6 months)
        if (request.Data.Fields.TryGetValue("customer_id", out var customerId))
        {
            var userId = Guid.Parse(customerId.StringValue);
            var sixMonthsAgo = DateTime.UtcNow.AddMonths(-6);
            var recentClaims = await _claimRepository.FindAsync(c => c.CustomerId == userId && c.CreatedAt >= sixMonthsAgo, cancellationToken);

            if (recentClaims.Count() >= 3)
            {
                response.IsFraudDetected = true;
                response.FraudScore += 50;
                response.TriggeredRules.Add("FR-094: Frequent claims (>3 in 6 months)");
            }
        }

        // FR-094: Rapid policy-to-claim (<48hrs)
        if (request.Data.Fields.TryGetValue("policy_id", out var policyId))
        {
            var policy = await _policyRepository.GetByIdAsync(Guid.Parse(policyId.StringValue), cancellationToken);
            if (policy != null && (DateTime.UtcNow - (policy.IssuedAt ?? DateTime.UtcNow)).TotalHours < 48)
            {
                response.IsFraudDetected = true;
                response.FraudScore += 60;
                response.TriggeredRules.Add("FR-094: Rapid policy-to-claim (<48hrs)");
            }
        }
    }
}
