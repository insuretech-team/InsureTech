using Microsoft.EntityFrameworkCore;

namespace PoliSync.Infrastructure.Persistence;

/// <summary>
/// Main EF Core DbContext for the InsureTech platform.
/// Module-specific entity configurations are applied via IEntityTypeConfiguration.
/// </summary>
public class PoliSyncDbContext : DbContext
{
    public PoliSyncDbContext(DbContextOptions<PoliSyncDbContext> options) : base(options) { }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // Apply all IEntityTypeConfiguration<T> from all loaded PoliSync assemblies
        // This allows Vertical Slices (Modules) to stay independent but still contribute to the schema
        var assemblies = AppDomain.CurrentDomain.GetAssemblies()
            .Where(a => a.FullName?.StartsWith("PoliSync.") == true);
        
        foreach (var assembly in assemblies)
        {
            modelBuilder.ApplyConfigurationsFromAssembly(assembly);
        }

        // Default schema
        modelBuilder.HasDefaultSchema("insurance_schema");
    }

    public override async Task<int> SaveChangesAsync(CancellationToken ct = default)
    {
        // Set audit timestamps
        foreach (var entry in ChangeTracker.Entries())
        {
            if (entry.State == EntityState.Added)
            {
                if (entry.Metadata.FindProperty("created_at") != null)
                    entry.Property("created_at").CurrentValue = DateTime.UtcNow;
                if (entry.Metadata.FindProperty("updated_at") != null)
                    entry.Property("updated_at").CurrentValue = DateTime.UtcNow;
            }
            else if (entry.State == EntityState.Modified)
            {
                if (entry.Metadata.FindProperty("updated_at") != null)
                    entry.Property("updated_at").CurrentValue = DateTime.UtcNow;
            }
        }

        return await base.SaveChangesAsync(ct);
    }
}
