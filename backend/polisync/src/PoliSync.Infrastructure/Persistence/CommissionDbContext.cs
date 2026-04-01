using Microsoft.EntityFrameworkCore;

namespace PoliSync.Infrastructure.Persistence;

/// <summary>
/// Separate DbContext for commission_schema — keeps commission data
/// isolated from insurance_schema for schema-level multi-tenancy.
/// </summary>
public sealed class CommissionDbContext : DbContext
{
    public CommissionDbContext(DbContextOptions<CommissionDbContext> options)
        : base(options)
    {
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);
        modelBuilder.HasDefaultSchema("commission_schema");
    }
}
