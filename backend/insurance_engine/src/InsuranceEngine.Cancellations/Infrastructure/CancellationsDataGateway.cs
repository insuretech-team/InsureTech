using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using InsuranceEngine.Cancellations;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Cancellations.Infrastructure;

public class CancellationsDbContext : DbContext
{
    public CancellationsDbContext(DbContextOptions<CancellationsDbContext> options) : base(options) { }

    public DbSet<PolicyEntity> Policies => Set<PolicyEntity>();
    public DbSet<RefundEntity> Refunds => Set<RefundEntity>();
    public DbSet<CancellationEntity> Cancellations => Set<CancellationEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");
        modelBuilder.Entity<PolicyEntity>().ToTable("policies");
        modelBuilder.Entity<RefundEntity>().ToTable("refunds");
        modelBuilder.Entity<CancellationEntity>().ToTable("cancellations");
    }
}

public class SqlCancellationDataGateway : ICancellationDataGateway
{
    private readonly CancellationsDbContext _context;
    private readonly ILogger<SqlCancellationDataGateway> _logger;

    public SqlCancellationDataGateway(CancellationsDbContext context, ILogger<SqlCancellationDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default)
    {
        var policyId = request.PolicyId;
        var policy = await _context.Policies.FirstOrDefaultAsync(p => p.PolicyId.ToString() == policyId, ct);

        if (policy == null)
        {
            return new CancelPolicyResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "NOT_FOUND", Message = "Policy not found" }
            };
        }

        policy.Status = "CANCELLED";
        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Cancelled policy {PolicyId}", policyId);

        return new CancelPolicyResponse { Message = "Policy cancelled" };
    }

    public async Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default)
    {
        var policyId = request.PolicyId;
        var policy = await _context.Policies.FirstOrDefaultAsync(p => p.PolicyId.ToString() == policyId, ct);

        if (policy == null)
        {
            return new ApproveCancellationResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "NOT_FOUND", Message = "Policy not found" }
            };
        }

        policy.Status = "CANCELLED";
        policy.UpdatedAt = DateTime.UtcNow;

        var refundId = Guid.NewGuid().ToString();
        var refund = new RefundEntity
        {
            RefundId = refundId,
            RefundNumber = $"REF-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}",
            PolicyId = policyId,
            RefundType = "CANCELLATION",
            RefundAmount = policy.PremiumAmount * 0.8m,
            RefundCurrency = policy.PremiumCurrency,
            Status = "PENDING",
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        _context.Refunds.Add(refund);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Approved cancellation for policy {PolicyId}, Refund: {RefundNumber}", policyId, refund.RefundNumber);

        return new ApproveCancellationResponse { Message = "Cancellation approved" };
    }
}
