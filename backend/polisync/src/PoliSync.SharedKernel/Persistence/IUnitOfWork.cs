namespace PoliSync.SharedKernel.Persistence;

/// <summary>
/// Unit of Work abstraction for transactional persistence.
/// </summary>
public interface IUnitOfWork
{
    Task<int> SaveChangesAsync(CancellationToken ct = default);
}
