using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using InsuranceEngine.Renewals;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Renewals.Infrastructure;

public class RenewalsDbContext : DbContext
{
    public RenewalsDbContext(DbContextOptions<RenewalsDbContext> options) : base(options) { }

    public DbSet<PolicyEntity> Policies => Set<PolicyEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");
        modelBuilder.Entity<PolicyEntity>().ToTable("policies");
    }
}

public class SqlRenewalDataGateway : IRenewalDataGateway
{
    private readonly RenewalsDbContext _context;
    private readonly ILogger<SqlRenewalDataGateway> _logger;

    public SqlRenewalDataGateway(RenewalsDbContext context, ILogger<SqlRenewalDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default)
    {
        var id = Guid.TryParse(request.PolicyId, out var pid) ? pid : Guid.Empty;
        var policy = await _context.Policies.FindAsync([id], ct);

        if (policy == null)
        {
            return new RenewPolicyTenureResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "NOT_FOUND", Message = "Policy not found" }
            };
        }

        var newStartDate = policy.EndDate;
        var newEndDate = newStartDate.AddMonths(request.TenureMonths);

        policy.TenureMonths = request.TenureMonths;
        policy.StartDate = newStartDate;
        policy.EndDate = newEndDate;
        policy.Status = "ACTIVE";
        policy.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Renewed policy {PolicyId} for {TenureMonths} months", request.PolicyId, request.TenureMonths);

        return new RenewPolicyTenureResponse { Message = "Policy renewed successfully" };
    }
}
