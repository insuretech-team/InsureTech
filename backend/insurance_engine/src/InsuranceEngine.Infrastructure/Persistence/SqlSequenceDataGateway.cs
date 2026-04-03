using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Infrastructure.DataGateways;

namespace InsuranceEngine.Infrastructure.Persistence;

/// <summary>
/// EF Core/SQL implementation for database sequence generation.
/// </summary>
public class SqlSequenceDataGateway : ISequenceDataGateway
{
    private readonly InsuranceDbContext _dbContext;

    public SqlSequenceDataGateway(InsuranceDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    public async Task<long> GetNextSequenceValueAsync(string sequenceName, CancellationToken ct = default)
    {
        if (_dbContext.Database.ProviderName == "Microsoft.EntityFrameworkCore.Sqlite")
        {
            // Fallback for Sqlite (which doesn't have native sequences in the same way)
            // This is a simplified mock for testing/development.
            // In a real Sqlite prod environment (rare for this scale), we'd use a table.
            return (long)(DateTime.UtcNow.Ticks % 1000000); 
        }

        var connection = _dbContext.Database.GetDbConnection();
        if (connection.State != System.Data.ConnectionState.Open)
            await connection.OpenAsync(ct);

        using var cmd = connection.CreateCommand();
        cmd.CommandText = $"SELECT nextval('{sequenceName}')";
        var result = await cmd.ExecuteScalarAsync(ct);
        return Convert.ToInt64(result);
    }
}
