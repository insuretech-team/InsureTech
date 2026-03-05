using PoliSync.SharedKernel.Persistence;

namespace PoliSync.Infrastructure.Persistence;

/// <summary>
/// Unit of Work wrapping PoliSyncDbContext.SaveChangesAsync.
/// </summary>
public class UnitOfWork : IUnitOfWork
{
    private readonly PoliSyncDbContext _db;

    public UnitOfWork(PoliSyncDbContext db) => _db = db;

    public Task<int> SaveChangesAsync(CancellationToken ct = default)
        => _db.SaveChangesAsync(ct);
}
