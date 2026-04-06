using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Underwriting.Infrastructure;

public class UnderwritingDbContext : DbContext
{
    public UnderwritingDbContext(DbContextOptions<UnderwritingDbContext> options) : base(options) { }

    public DbSet<QuoteEntity> Quotes => Set<QuoteEntity>();
    public DbSet<UnderwritingDecisionEntity> UnderwritingDecisions => Set<UnderwritingDecisionEntity>();
    public DbSet<HealthDeclarationEntity> HealthDeclarations => Set<HealthDeclarationEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");
        modelBuilder.Entity<QuoteEntity>().ToTable("quotes");
        modelBuilder.Entity<UnderwritingDecisionEntity>().ToTable("underwriting_decisions");
        modelBuilder.Entity<HealthDeclarationEntity>().ToTable("health_declarations");
    }
}
