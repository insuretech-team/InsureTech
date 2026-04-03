namespace InsuranceEngine.SharedKernel.Infrastructure.DataGateways;

/// <summary>
/// Abstraction for database sequence generation to avoid direct DbContext dependency in Handlers.
/// </summary>
public interface ISequenceDataGateway
{
    Task<long> GetNextSequenceValueAsync(string sequenceName, CancellationToken ct = default);
}
