using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Common.V1;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;

namespace InsuranceEngine.Policy.Infrastructure;

/// <summary>
/// In-memory implementation of IPolicyDataGateway for testing and development.
/// No external dependencies (PostgreSQL, Go backend) required.
/// </summary>
public class InMemoryPolicyDataGateway : IPolicyDataGateway
{
    private readonly Dictionary<string, PolicyEntity> _policies = new();
    private int _policyCounter = 1;
    private readonly ILogger<InMemoryPolicyDataGateway> _logger;

    public InMemoryPolicyDataGateway(ILogger<InMemoryPolicyDataGateway> logger)
    {
        _logger = logger;
    }

    public Task<CreatePolicyResponse> CreatePolicyAsync(CreatePolicyRequest request, CancellationToken ct = default)
    {
        var policyId = Guid.NewGuid().ToString();
        var policyNumber = $"POL-{DateTime.UtcNow.Year}-{_policyCounter++:D5}";
        var now = DateTime.UtcNow;

        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = policyNumber,
            ProductId = request.ProductId,
            CustomerId = request.CustomerId,
            Status = PolicyStatus.PendingPayment,
            PremiumAmount = request.PremiumAmount ?? new Money(),
            SumInsured = request.SumInsured ?? new Money(),
            TenureMonths = request.TenureMonths > 0 ? request.TenureMonths : 12,
            StartDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(now),
            EndDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(now.AddMonths(request.TenureMonths > 0 ? request.TenureMonths : 12))
        };

        if (request.Nominees?.Count > 0)
        {
            foreach (var n in request.Nominees)
            {
                policy.Nominees.Add(n);
            }
        }

        _policies[policyId] = policy;
        _logger.LogInformation("InMemory: Created policy {PolicyNumber}", policyNumber);

        return Task.FromResult(new CreatePolicyResponse
        {
            PolicyId = policyId,
            PolicyNumber = policyNumber
        });
    }

    public Task<GetPolicyResponse> GetPolicyAsync(string policyId, CancellationToken ct = default)
    {
        if (_policies.TryGetValue(policyId, out var policy))
        {
            return Task.FromResult(new GetPolicyResponse { Policy = policy });
        }
        return Task.FromResult(new GetPolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } });
    }

    public Task<ListUserPoliciesResponse> ListUserPoliciesAsync(ListUserPoliciesRequest request, CancellationToken ct = default)
    {
        var query = _policies.Values.AsQueryable();

        if (!string.IsNullOrEmpty(request.CustomerId))
            query = query.Where(p => p.CustomerId == request.CustomerId);

        var page = request.Page > 0 ? request.Page : 1;
        var pageSize = request.PageSize > 0 ? request.PageSize : 10;

        var policies = query.Skip((page - 1) * pageSize).Take(pageSize).ToList();

        var response = new ListUserPoliciesResponse();
        response.Policies.AddRange(policies);

        return Task.FromResult(response);
    }

    public Task<UpdatePolicyResponse> UpdatePolicyAsync(string policyId, List<Nominee>? nominees, string? address, CancellationToken ct = default)
    {
        if (!_policies.TryGetValue(policyId, out var policy))
        {
            return Task.FromResult(new UpdatePolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } });
        }

        if (nominees != null)
        {
            policy.Nominees.Clear();
            foreach (var n in nominees)
            {
                policy.Nominees.Add(n);
            }
        }

        _logger.LogInformation("InMemory: Updated policy {PolicyId}", policyId);
        return Task.FromResult(new UpdatePolicyResponse());
    }

    public Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default)
    {
        if (!_policies.TryGetValue(request.PolicyId, out var policy))
        {
            return Task.FromResult(new CancelPolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } });
        }

        policy.Status = PolicyStatus.Cancelled;
        _logger.LogInformation("InMemory: Cancelled policy {PolicyId}", request.PolicyId);

        return Task.FromResult(new CancelPolicyResponse());
    }

    public Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default)
    {
        if (!_policies.TryGetValue(request.PolicyId, out var policy))
        {
            return Task.FromResult(new RenewPolicyTenureResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } });
        }

        policy.TenureMonths = request.TenureMonths;
        _logger.LogInformation("InMemory: Renewed policy {PolicyId}", request.PolicyId);

        return Task.FromResult(new RenewPolicyTenureResponse());
    }

    public Task<GeneratePolicyDocumentResponse> GeneratePolicyDocumentAsync(string policyId, CancellationToken ct = default)
    {
        if (!_policies.ContainsKey(policyId))
        {
            return Task.FromResult(new GeneratePolicyDocumentResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } });
        }

        var docUrl = $"https://storage.insuretech/policies/{policyId}.pdf";
        _logger.LogInformation("InMemory: Generated document for policy {PolicyId}", policyId);

        return Task.FromResult(new GeneratePolicyDocumentResponse { DocumentUrl = docUrl });
    }

    public Task<IssuePolicyResponse> IssuePolicyAsync(IssuePolicyRequest request, CancellationToken ct = default)
    {
        if (!_policies.TryGetValue(request.PolicyId, out var policy))
        {
            return Task.FromResult(new IssuePolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } });
        }

        policy.Status = PolicyStatus.Active;
        _logger.LogInformation("InMemory: Issued policy {PolicyId}", request.PolicyId);

        return Task.FromResult(new IssuePolicyResponse());
    }

    public Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default)
    {
        if (!_policies.TryGetValue(request.PolicyId, out var policy))
        {
            return Task.FromResult(new ApproveCancellationResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } });
        }

        policy.Status = PolicyStatus.Cancelled;
        _logger.LogInformation("InMemory: Approved cancellation for policy {PolicyId}", request.PolicyId);

        return Task.FromResult(new ApproveCancellationResponse());
    }
}
