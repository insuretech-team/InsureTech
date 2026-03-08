using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Policy.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Policy.Infrastructure;

namespace PoliSync.Policy.GrpcServices;

public sealed class PolicyGrpcService : PolicyService.PolicyServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<PolicyGrpcService> _logger;
    private readonly IPolicyDataGateway _policyDataGateway;

    public PolicyGrpcService(
        IMediator mediator,
        ILogger<PolicyGrpcService> logger,
        IPolicyDataGateway policyDataGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _policyDataGateway = policyDataGateway;
    }

    public override async Task<CreatePolicyResponse> CreatePolicy(CreatePolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.ProductId) || string.IsNullOrWhiteSpace(request.CustomerId))
        {
            return new CreatePolicyResponse
            {
                Error = BuildError("VALIDATION_ERROR", "ProductId and CustomerId are required")
            };
        }

        try
        {
            var now = DateTime.UtcNow;
            var tenureMonths = request.TenureMonths <= 0 ? 12 : request.TenureMonths;
            var premium = NormalizeMoney(request.PremiumAmount, 120_000);
            var vat = NewMoney((long)Math.Round(premium.Amount * 0.15));
            var fee = NewMoney(2_000);

            var policy = new Insuretech.Policy.Entity.V1.Policy
            {
                PolicyId = Guid.NewGuid().ToString("N"),
                PolicyNumber = BuildPolicyNumber(),
                ProductId = request.ProductId,
                CustomerId = request.CustomerId,
                PartnerId = request.PartnerId,
                AgentId = request.AgentId,
                Status = PolicyStatus.PendingPayment,
                PremiumAmount = premium,
                SumInsured = NormalizeMoney(request.SumInsured, 1_000_000),
                TenureMonths = tenureMonths,
                StartDate = Timestamp.FromDateTime(now.Date),
                EndDate = Timestamp.FromDateTime(now.Date.AddMonths(tenureMonths)),
                CreatedAt = Timestamp.FromDateTime(now),
                UpdatedAt = Timestamp.FromDateTime(now),
                PaymentFrequency = "MONTHLY",
                VatTax = vat,
                ServiceFee = fee,
                TotalPayable = NewMoney(premium.Amount + vat.Amount + fee.Amount),
                PremiumCurrency = premium.Currency,
                SumInsuredCurrency = request.SumInsured?.Currency ?? "BDT",
                ProposerDetails = request.Applicant ?? new Applicant()
            };

            policy.Nominees.AddRange(request.Nominees);
            policy.Riders.AddRange(request.Riders);

            var created = await _policyDataGateway.CreatePolicyAsync(policy, GetCancellationToken(context));
            _logger.LogInformation("Policy created via Go insurance service: {PolicyId}", created.PolicyId);

            return new CreatePolicyResponse
            {
                PolicyId = created.PolicyId,
                PolicyNumber = created.PolicyNumber,
                Message = "Policy created successfully"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "CreatePolicy failed with gRPC status {Status}", ex.StatusCode);
            return new CreatePolicyResponse
            {
                Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}")
            };
        }
    }

    public override async Task<GetPolicyResponse> GetPolicy(GetPolicyRequest request, ServerCallContext context)
    {
        try
        {
            var policy = await _policyDataGateway.GetPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (policy is null)
            {
                return new GetPolicyResponse { Error = BuildError("NOT_FOUND", "Policy not found") };
            }

            return new GetPolicyResponse { Policy = policy };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "GetPolicy failed with gRPC status {Status}", ex.StatusCode);
            return new GetPolicyResponse { Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}") };
        }
    }

    public override async Task<ListUserPoliciesResponse> ListUserPolicies(ListUserPoliciesRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;

        try
        {
            var policies = await _policyDataGateway.ListPoliciesAsync(request.CustomerId, page, pageSize, GetCancellationToken(context));

            var filtered = request.Status == PolicyStatus.Unspecified
                ? policies
                : policies.Where(x => x.Status == request.Status).ToList();

            var response = new ListUserPoliciesResponse
            {
                TotalCount = filtered.Count
            };
            response.Policies.AddRange(filtered);
            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "ListUserPolicies failed with gRPC status {Status}", ex.StatusCode);
            return new ListUserPoliciesResponse { Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}") };
        }
    }

    public override async Task<UpdatePolicyResponse> UpdatePolicy(UpdatePolicyRequest request, ServerCallContext context)
    {
        try
        {
            var policy = await _policyDataGateway.GetPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (policy is null)
            {
                return new UpdatePolicyResponse { Error = BuildError("NOT_FOUND", "Policy not found") };
            }

            if (policy.Status is PolicyStatus.Cancelled or PolicyStatus.Expired)
            {
                return new UpdatePolicyResponse { Error = BuildError("INVALID_STATE", "Cancelled or expired policy cannot be updated") };
            }

            policy.Nominees.Clear();
            policy.Nominees.AddRange(request.Nominees);

            if (!string.IsNullOrWhiteSpace(request.Address))
            {
                var proposer = policy.ProposerDetails ?? new Applicant();
                proposer.Address = request.Address;
                policy.ProposerDetails = proposer;
            }

            policy.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await _policyDataGateway.UpdatePolicyAsync(policy, GetCancellationToken(context));

            return new UpdatePolicyResponse { Message = "Policy updated" };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "UpdatePolicy failed with gRPC status {Status}", ex.StatusCode);
            return new UpdatePolicyResponse { Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}") };
        }
    }

    public override async Task<CancelPolicyResponse> CancelPolicy(CancelPolicyRequest request, ServerCallContext context)
    {
        try
        {
            var policy = await _policyDataGateway.GetPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (policy is null)
            {
                return new CancelPolicyResponse { Error = BuildError("NOT_FOUND", "Policy not found") };
            }

            if (policy.Status == PolicyStatus.Cancelled)
            {
                return new CancelPolicyResponse { Error = BuildError("INVALID_STATE", "Policy already cancelled") };
            }

            policy.Status = PolicyStatus.Cancelled;
            policy.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await _policyDataGateway.UpdatePolicyAsync(policy, GetCancellationToken(context));

            var refund = NewMoney((long)Math.Round(policy.PremiumAmount.Amount * 0.60), policy.PremiumAmount.Currency);
            return new CancelPolicyResponse
            {
                Message = "Policy cancelled",
                RefundAmount = refund
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "CancelPolicy failed with gRPC status {Status}", ex.StatusCode);
            return new CancelPolicyResponse { Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}") };
        }
    }

    public override async Task<RenewPolicyResponse> RenewPolicy(RenewPolicyRequest request, ServerCallContext context)
    {
        try
        {
            var current = await _policyDataGateway.GetPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (current is null)
            {
                return new RenewPolicyResponse { Error = BuildError("NOT_FOUND", "Policy not found") };
            }

            var tenureMonths = request.TenureMonths <= 0 ? current.TenureMonths : request.TenureMonths;
            var startDate = SafeDate(current.EndDate, DateTime.UtcNow.Date);
            var now = DateTime.UtcNow;

            var renewed = ClonePolicy(current);
            renewed.PolicyId = Guid.NewGuid().ToString("N");
            renewed.PolicyNumber = BuildPolicyNumber();
            renewed.Status = PolicyStatus.PendingPayment;
            renewed.StartDate = Timestamp.FromDateTime(startDate);
            renewed.EndDate = Timestamp.FromDateTime(startDate.AddMonths(tenureMonths));
            renewed.TenureMonths = tenureMonths;
            renewed.CreatedAt = Timestamp.FromDateTime(now);
            renewed.UpdatedAt = Timestamp.FromDateTime(now);
            renewed.IssuedAt = new Timestamp();
            renewed.PaymentGatewayReference = string.Empty;
            renewed.ReceiptNumber = string.Empty;
            renewed.PolicyDocumentUrl = string.Empty;

            if (request.UpdateNominees)
            {
                renewed.Nominees.Clear();
                renewed.Nominees.AddRange(request.Nominees);
            }

            var created = await _policyDataGateway.CreatePolicyAsync(renewed, GetCancellationToken(context));
            return new RenewPolicyResponse
            {
                NewPolicyId = created.PolicyId,
                NewPolicyNumber = created.PolicyNumber,
                PremiumAmount = created.PremiumAmount,
                Message = "Policy renewed"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "RenewPolicy failed with gRPC status {Status}", ex.StatusCode);
            return new RenewPolicyResponse { Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}") };
        }
    }

    public override async Task<GeneratePolicyDocumentResponse> GeneratePolicyDocument(GeneratePolicyDocumentRequest request, ServerCallContext context)
    {
        try
        {
            var policy = await _policyDataGateway.GetPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (policy is null)
            {
                return new GeneratePolicyDocumentResponse { Error = BuildError("NOT_FOUND", "Policy not found") };
            }

            var url = $"https://docs.polisync.local/policies/{policy.PolicyId}.pdf";
            policy.PolicyDocumentUrl = url;
            policy.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await _policyDataGateway.UpdatePolicyAsync(policy, GetCancellationToken(context));

            return new GeneratePolicyDocumentResponse
            {
                DocumentUrl = url,
                QrCode = $"QR-{policy.PolicyId[..Math.Min(12, policy.PolicyId.Length)].ToUpperInvariant()}"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "GeneratePolicyDocument failed with gRPC status {Status}", ex.StatusCode);
            return new GeneratePolicyDocumentResponse { Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}") };
        }
    }

    public override async Task<IssuePolicyResponse> IssuePolicy(IssuePolicyRequest request, ServerCallContext context)
    {
        try
        {
            var policy = await _policyDataGateway.GetPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (policy is null)
            {
                return new IssuePolicyResponse { Error = BuildError("NOT_FOUND", "Policy not found") };
            }

            policy.Status = PolicyStatus.Active;
            policy.IssuedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            policy.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            policy.QuoteId = request.QuoteId;
            policy.PaymentGatewayReference = request.PaymentId;
            policy.ReceiptNumber = $"RCPT-{DateTime.UtcNow:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}";

            var updated = await _policyDataGateway.UpdatePolicyAsync(policy, GetCancellationToken(context));
            return new IssuePolicyResponse
            {
                Policy = updated,
                Message = "Policy issued"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "IssuePolicy failed with gRPC status {Status}", ex.StatusCode);
            return new IssuePolicyResponse { Error = BuildError("UPSTREAM_ERROR", $"Insurance service error: {ex.StatusCode}") };
        }
    }

    private static DateTime SafeDate(Timestamp? timestamp, DateTime fallbackUtcDate)
    {
        if (timestamp is null || timestamp.Seconds <= 0)
        {
            return DateTime.SpecifyKind(fallbackUtcDate, DateTimeKind.Utc);
        }

        return timestamp.ToDateTime().Date;
    }

    private static string BuildPolicyNumber()
        => $"LP-{DateTime.UtcNow:yyyy}-{Random.Shared.Next(100000, 999999)}";

    private static Money NormalizeMoney(Money? source, long fallbackAmount)
        => new()
        {
            Amount = source?.Amount > 0 ? source.Amount : fallbackAmount,
            Currency = string.IsNullOrWhiteSpace(source?.Currency) ? "BDT" : source.Currency
        };

    private static Money NewMoney(long amount, string currency = "BDT")
        => new() { Amount = amount, Currency = currency };

    private static Error BuildError(string code, string message)
        => new() { Code = code, Message = message };

    private static CancellationToken GetCancellationToken(ServerCallContext? context)
        => context?.CancellationToken ?? CancellationToken.None;

    private static Insuretech.Policy.Entity.V1.Policy ClonePolicy(Insuretech.Policy.Entity.V1.Policy source)
    {
        var clone = new Insuretech.Policy.Entity.V1.Policy
        {
            ProductId = source.ProductId,
            CustomerId = source.CustomerId,
            PartnerId = source.PartnerId,
            AgentId = source.AgentId,
            QuoteId = source.QuoteId,
            UnderwritingDecisionId = source.UnderwritingDecisionId,
            Status = source.Status,
            PremiumAmount = source.PremiumAmount,
            SumInsured = source.SumInsured,
            TenureMonths = source.TenureMonths,
            StartDate = source.StartDate,
            EndDate = source.EndDate,
            IssuedAt = source.IssuedAt,
            CreatedAt = source.CreatedAt,
            UpdatedAt = source.UpdatedAt,
            PolicyDocumentUrl = source.PolicyDocumentUrl,
            DeletedAt = source.DeletedAt,
            PaymentFrequency = source.PaymentFrequency,
            VatTax = source.VatTax,
            ServiceFee = source.ServiceFee,
            TotalPayable = source.TotalPayable,
            PaymentGatewayReference = source.PaymentGatewayReference,
            ReceiptNumber = source.ReceiptNumber,
            ProposerDetails = source.ProposerDetails,
            OccupationRiskClass = source.OccupationRiskClass,
            HasExistingPolicies = source.HasExistingPolicies,
            ClaimsHistorySummary = source.ClaimsHistorySummary,
            ProviderName = source.ProviderName,
            EnrollmentStartDate = source.EnrollmentStartDate,
            EnrollmentEndDate = source.EnrollmentEndDate,
            UnderwritingData = source.UnderwritingData,
            PremiumCurrency = source.PremiumCurrency,
            SumInsuredCurrency = source.SumInsuredCurrency
        };
        clone.Nominees.AddRange(source.Nominees);
        clone.Riders.AddRange(source.Riders);
        return clone;
    }
}
