using Microsoft.Extensions.Logging;
using Microsoft.EntityFrameworkCore;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Common.V1;
using PolicyEntity = InsuranceEngine.SharedKernel.Persistence.Entities.PolicyEntity;
using NomineeEntity = InsuranceEngine.SharedKernel.Persistence.Entities.PolicyNomineeEntity;
using ProtoPolicy = Insuretech.Policy.Entity.V1.Policy;
using ProtoNominee = Insuretech.Policy.Entity.V1.Nominee;

namespace InsuranceEngine.Policy.Infrastructure;

public class SqlPolicyDataGateway : IPolicyDataGateway
{
    private readonly PolicyDbContext _context;
    private readonly ILogger<SqlPolicyDataGateway> _logger;

    public SqlPolicyDataGateway(PolicyDbContext context, ILogger<SqlPolicyDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<CreatePolicyResponse> CreatePolicyAsync(CreatePolicyRequest request, CancellationToken ct = default)
    {
        var policyId = Guid.NewGuid();
        var policyNumber = $"POL-{DateTime.UtcNow.Year}-{DateTime.UtcNow:MMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        var now = DateTime.UtcNow;

        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = policyNumber,
            ProductId = Guid.TryParse(request.ProductId, out var pid) ? pid : Guid.Empty,
            CustomerId = Guid.TryParse(request.CustomerId, out var cid) ? cid : Guid.Empty,
            Status = "PENDING_PAYMENT",
            PremiumAmount = request.PremiumAmount?.Amount ?? 0,
            PremiumCurrency = request.PremiumAmount?.Currency ?? "BDT",
            SumInsuredAmount = request.SumInsured?.Amount ?? 0,
            SumInsuredCurrency = request.SumInsured?.Currency ?? "BDT",
            TenureMonths = request.TenureMonths > 0 ? request.TenureMonths : 12,
            StartDate = now,
            EndDate = now.AddMonths(request.TenureMonths > 0 ? request.TenureMonths : 12),
            CreatedAt = now,
            UpdatedAt = now
        };

        if (request.Nominees?.Count > 0)
        {
            foreach (var n in request.Nominees)
            {
                var nominee = new NomineeEntity
                {
                    NomineeId = Guid.NewGuid(),
                    PolicyId = policyId,
                    FullName = n.FullName ?? "",
                    Relationship = n.Relationship ?? "",
                    SharePercentage = n.SharePercentage > 0 ? n.SharePercentage : 100,
                    DateOfBirth = n.NomineeDobText != null ? DateTime.TryParse(n.NomineeDobText, out var dob) ? dob : DateTime.UnixEpoch : DateTime.UnixEpoch,
                    NidNumber = n.NidNumber,
                    PhoneNumber = n.PhoneNumber,
                    NomineeDobText = n.NomineeDobText,
                    CreatedAt = now,
                    UpdatedAt = now
                };
                policy.Nominees.Add(nominee);
            }
        }

        _context.Policies.Add(policy);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Created policy {PolicyNumber}", policyNumber);

        return new CreatePolicyResponse
        {
            PolicyId = policyId.ToString(),
            PolicyNumber = policyNumber
        };
    }

    public async Task<GetPolicyResponse> GetPolicyAsync(string policyId, CancellationToken ct = default)
    {
        var id = Guid.TryParse(policyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies
            .Include(p => p.Nominees)
            .FirstOrDefaultAsync(p => p.PolicyId == id, ct);

        if (policy == null)
        {
            return new GetPolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
        }

        return new GetPolicyResponse { Policy = MapToProto(policy) };
    }

    public async Task<ListUserPoliciesResponse> ListUserPoliciesAsync(ListUserPoliciesRequest request, CancellationToken ct = default)
    {
        var query = _context.Policies.AsQueryable();

        if (!string.IsNullOrEmpty(request.CustomerId) && Guid.TryParse(request.CustomerId, out var cid))
        {
            query = query.Where(p => p.CustomerId == cid);
        }

        var page = request.Page > 0 ? request.Page : 1;
        var pageSize = request.PageSize > 0 ? request.PageSize : 10;

        var policies = await query
            .OrderByDescending(p => p.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .Include(p => p.Nominees)
            .ToListAsync(ct);

        var response = new ListUserPoliciesResponse();
        response.Policies.AddRange(policies.Select(MapToProto));

        return response;
    }

    public async Task<UpdatePolicyResponse> UpdatePolicyAsync(string policyId, List<Nominee>? nominees, string? address, CancellationToken ct = default)
    {
        var id = Guid.TryParse(policyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies
            .Include(p => p.Nominees)
            .FirstOrDefaultAsync(p => p.PolicyId == id, ct);

        if (policy == null)
        {
            return new UpdatePolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
        }

        if (nominees != null)
        {
            _context.PolicyNominees.RemoveRange(policy.Nominees);
            var now = DateTime.UtcNow;
            foreach (var n in nominees)
            {
                var nominee = new NomineeEntity
                {
                    NomineeId = Guid.NewGuid(),
                    PolicyId = id,
                    FullName = n.FullName ?? "",
                    Relationship = n.Relationship ?? "",
                    SharePercentage = n.SharePercentage > 0 ? n.SharePercentage : 100,
                    NidNumber = n.NidNumber,
                    PhoneNumber = n.PhoneNumber,
                    NomineeDobText = n.NomineeDobText,
                    CreatedAt = now,
                    UpdatedAt = now
                };
                _context.PolicyNominees.Add(nominee);
            }
        }

        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Updated policy {PolicyId}", policyId);
        return new UpdatePolicyResponse();
    }

    public async Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default)
    {
        var id = Guid.TryParse(request.PolicyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies.FindAsync([id], ct);

        if (policy == null)
        {
            return new CancelPolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
        }

        policy.Status = "CANCELLED";
        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Cancelled policy {PolicyId}", request.PolicyId);
        return new CancelPolicyResponse();
    }

    public async Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default)
    {
        var id = Guid.TryParse(request.PolicyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies.FindAsync([id], ct);

        if (policy == null)
        {
            return new RenewPolicyTenureResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
        }

        policy.TenureMonths = request.TenureMonths;
        policy.StartDate = policy.EndDate;
        policy.EndDate = policy.EndDate.AddMonths(request.TenureMonths);
        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Renewed policy {PolicyId}", request.PolicyId);
        return new RenewPolicyTenureResponse();
    }

    public async Task<GeneratePolicyDocumentResponse> GeneratePolicyDocumentAsync(string policyId, CancellationToken ct = default)
    {
        var id = Guid.TryParse(policyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies.FindAsync([id], ct);

        if (policy == null)
        {
            return new GeneratePolicyDocumentResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
        }

        var docUrl = $"https://storage.insuretech/policies/{policy.PolicyNumber}.pdf";
        policy.PolicyDocumentUrl = docUrl;
        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Generated document for policy {PolicyId}", policyId);
        return new GeneratePolicyDocumentResponse { DocumentUrl = docUrl };
    }

    public async Task<IssuePolicyResponse> IssuePolicyAsync(IssuePolicyRequest request, CancellationToken ct = default)
    {
        var id = Guid.TryParse(request.PolicyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies.FindAsync([id], ct);

        if (policy == null)
        {
            return new IssuePolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
        }

        policy.Status = "ACTIVE";
        policy.IssuedAt = DateTime.UtcNow;
        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Issued policy {PolicyId}", request.PolicyId);
        return new IssuePolicyResponse();
    }

    public async Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default)
    {
        var id = Guid.TryParse(request.PolicyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies.FindAsync([id], ct);

        if (policy == null)
        {
            return new ApproveCancellationResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
        }

        policy.Status = "CANCELLED";
        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Approved cancellation for policy {PolicyId}", request.PolicyId);
        return new ApproveCancellationResponse();
    }

    private static ProtoPolicy MapToProto(PolicyEntity entity)
    {
        var proto = new ProtoPolicy
        {
            PolicyId = entity.PolicyId.ToString(),
            PolicyNumber = entity.PolicyNumber,
            ProductId = entity.ProductId.ToString(),
            CustomerId = entity.CustomerId.ToString(),
            Status = Enum.Parse<PolicyStatus>(entity.Status.Replace("_", ""), true),
            PremiumAmount = new Money { Amount = entity.PremiumAmount, Currency = entity.PremiumCurrency },
            SumInsured = new Money { Amount = entity.SumInsuredAmount, Currency = entity.SumInsuredCurrency },
            TenureMonths = entity.TenureMonths,
            StartDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.StartDate),
            EndDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.EndDate)
        };

        if (entity.Nominees != null)
        {
            foreach (var n in entity.Nominees)
            {
                proto.Nominees.Add(new ProtoNominee
                {
                    FullName = n.FullName,
                    Relationship = n.Relationship,
                    SharePercentage = n.SharePercentage,
                    NidNumber = n.NidNumber ?? "",
                    PhoneNumber = n.PhoneNumber ?? "",
                    NomineeDobText = n.NomineeDobText ?? ""
                });
            }
        }

        return proto;
    }
}
