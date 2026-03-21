using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;

namespace InsuranceEngine.Policy.Domain.Services;

public class EndorsementNumberGenerator
{
    private readonly IEndorsementRepository _repo;

    public EndorsementNumberGenerator(IEndorsementRepository repo)
    {
        _repo = repo;
    }

    public async Task<string> GenerateNumberAsync(string policyNumber)
    {
        var seq = await _repo.GetNextSequenceNumberAsync();
        // Format: POL-001/END-01
        return $"{policyNumber}/END-{seq:D2}";
    }
}
