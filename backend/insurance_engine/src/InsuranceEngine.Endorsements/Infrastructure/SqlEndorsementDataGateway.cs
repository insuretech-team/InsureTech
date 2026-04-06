using InsuranceEngine.SharedKernel.Persistence.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Endorsements.Infrastructure;

public interface ISqlEndorsementDataGateway
{
    Task<EndorsementEntity?> GetByIdAsync(string endorsementId, CancellationToken ct = default);
    Task<EndorsementEntity?> GetByPolicyIdAsync(string policyId, CancellationToken ct = default);
    Task<List<EndorsementEntity>> GetByPolicyIdListAsync(string policyId, CancellationToken ct = default);
    Task<EndorsementEntity> CreateAsync(EndorsementEntity endorsement, CancellationToken ct = default);
    Task<EndorsementEntity> UpdateAsync(EndorsementEntity endorsement, CancellationToken ct = default);
    Task<EndorsementDocumentEntity?> GetDocumentByIdAsync(string documentId, CancellationToken ct = default);
    Task<EndorsementDocumentEntity> CreateDocumentAsync(EndorsementDocumentEntity document, CancellationToken ct = default);
}

public class SqlEndorsementDataGateway : ISqlEndorsementDataGateway
{
    private readonly EndorsementsDbContext _context;
    private readonly ILogger<SqlEndorsementDataGateway> _logger;

    public SqlEndorsementDataGateway(EndorsementsDbContext context, ILogger<SqlEndorsementDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<EndorsementEntity?> GetByIdAsync(string endorsementId, CancellationToken ct = default)
    {
        return await _context.Endorsements
            .FirstOrDefaultAsync(e => e.EndorsementId == endorsementId, ct);
    }

    public async Task<EndorsementEntity?> GetByPolicyIdAsync(string policyId, CancellationToken ct = default)
    {
        return await _context.Endorsements
            .Where(e => e.PolicyId == policyId)
            .OrderByDescending(e => e.CreatedAt)
            .FirstOrDefaultAsync(ct);
    }

    public async Task<List<EndorsementEntity>> GetByPolicyIdListAsync(string policyId, CancellationToken ct = default)
    {
        return await _context.Endorsements
            .Where(e => e.PolicyId == policyId)
            .OrderByDescending(e => e.CreatedAt)
            .ToListAsync(ct);
    }

    public async Task<EndorsementEntity> CreateAsync(EndorsementEntity endorsement, CancellationToken ct = default)
    {
        _context.Endorsements.Add(endorsement);
        await _context.SaveChangesAsync(ct);
        _logger.LogInformation("Created endorsement {EndorsementId}", endorsement.EndorsementId);
        return endorsement;
    }

    public async Task<EndorsementEntity> UpdateAsync(EndorsementEntity endorsement, CancellationToken ct = default)
    {
        endorsement.UpdatedAt = DateTime.UtcNow;
        _context.Endorsements.Update(endorsement);
        await _context.SaveChangesAsync(ct);
        _logger.LogInformation("Updated endorsement {EndorsementId}", endorsement.EndorsementId);
        return endorsement;
    }

    public async Task<EndorsementDocumentEntity?> GetDocumentByIdAsync(string documentId, CancellationToken ct = default)
    {
        return await _context.EndorsementDocuments
            .FirstOrDefaultAsync(d => d.DocumentId == documentId, ct);
    }

    public async Task<EndorsementDocumentEntity> CreateDocumentAsync(EndorsementDocumentEntity document, CancellationToken ct = default)
    {
        _context.EndorsementDocuments.Add(document);
        await _context.SaveChangesAsync(ct);
        _logger.LogInformation("Created endorsement document {DocumentId}", document.DocumentId);
        return document;
    }
}
