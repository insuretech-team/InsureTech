using System.Collections.Concurrent;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Crm.Entity.V1;

namespace PoliSync.CRM.Services;

public class InMemoryLeadRepository : ILeadRepository
{
    private readonly ConcurrentDictionary<string, Lead> _leads = new();
    private readonly ILogger<InMemoryLeadRepository> _logger;

    public InMemoryLeadRepository(ILogger<InMemoryLeadRepository> logger)
    {
        _logger = logger;
        SeedDefaultLeads();
    }

    public Task<Lead?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _leads.TryGetValue(id, out var lead);
        return Task.FromResult(lead);
    }

    public Task<IEnumerable<Lead>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_leads.Values.AsEnumerable());
    }

    public Task<IEnumerable<Lead>> GetByFilterAsync(LeadStatus? status, LeadSource? source, string? assignedAgentId, CancellationToken cancellationToken = default)
    {
        var query = _leads.Values.AsEnumerable();
        
        if (status.HasValue)
        {
            query = query.Where(l => l.LeadStatus == status.Value);
        }
        
        if (source.HasValue)
        {
            query = query.Where(l => l.LeadSource == source.Value);
        }
        
        if (!string.IsNullOrEmpty(assignedAgentId))
        {
            query = query.Where(l => l.AssignedAgentId == assignedAgentId);
        }
        
        return Task.FromResult<IEnumerable<Lead>>(query.OrderByDescending(l => l.CreatedAt));
    }

    public Task<IEnumerable<Lead>> GetByAssignedAgentAsync(string agentId, CancellationToken cancellationToken = default)
    {
        var leads = _leads.Values
            .Where(l => l.AssignedAgentId == agentId)
            .OrderByDescending(l => l.CreatedAt)
            .AsEnumerable();
        return Task.FromResult(leads);
    }

    public Task<Lead> CreateAsync(Lead lead, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(lead.LeadId))
        {
            lead.LeadId = Guid.NewGuid().ToString();
        }
        
        lead.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        lead.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        _leads[lead.LeadId] = lead;
        _logger.LogInformation("Created lead: {LeadId} - {FirstName} {LastName}", 
            lead.LeadId, lead.FirstName, lead.LastName);
        
        return Task.FromResult(lead);
    }

    public Task<Lead?> UpdateAsync(Lead lead, CancellationToken cancellationToken = default)
    {
        if (!_leads.ContainsKey(lead.LeadId))
        {
            return Task.FromResult<Lead?>(null);
        }

        lead.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _leads[lead.LeadId] = lead;
        
        _logger.LogInformation("Updated lead: {LeadId}", lead.LeadId);
        
        return Task.FromResult<Lead?>(lead);
    }

    public Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default)
    {
        if (permanent)
        {
            var result = _leads.TryRemove(id, out _);
            if (result)
            {
                _logger.LogInformation("Permanently deleted lead: {LeadId}", id);
            }
            return Task.FromResult(result);
        }
        else
        {
            if (_leads.TryGetValue(id, out var lead))
            {
                lead.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                _logger.LogInformation("Soft deleted lead: {LeadId}", id);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }

    private void SeedDefaultLeads()
    {
        var sampleLeads = new[]
        {
            new Lead
            {
                LeadId = Guid.NewGuid().ToString(),
                Title = "Mr",
                FirstName = "Rahim",
                LastName = "Khan",
                Gender = "Male",
                PhoneNumber = "+8801712345678",
                EmailAddress = "rahim.khan@email.com",
                Address = "Dhaka, Bangladesh",
                LeadSource = LeadSource.Website,
                LeadStatus = LeadStatus.New,
                LeadPriority = LeadPriority.High,
                LeadScore = 75,
                DesiredInsuranceType = "Term Life",
                DesiredCoverageAmount = 50000000, // 500,000 BDT
                QualificationStatus = QualificationStatus.InProgress,
                CreatedBy = "system"
            },
            new Lead
            {
                LeadId = Guid.NewGuid().ToString(),
                Title = "Ms",
                FirstName = "Fatima",
                LastName = "Ahmed",
                Gender = "Female",
                PhoneNumber = "+8801812345678",
                EmailAddress = "fatima.ahmed@email.com",
                Address = "Chittagong, Bangladesh",
                LeadSource = LeadSource.Referral,
                LeadStatus = LeadStatus.Contacted,
                LeadPriority = LeadPriority.Medium,
                LeadScore = 60,
                DesiredInsuranceType = "Health",
                DesiredCoverageAmount = 30000000, // 300,000 BDT
                QualificationStatus = QualificationStatus.Qualified,
                CreatedBy = "system"
            },
            new Lead
            {
                LeadId = Guid.NewGuid().ToString(),
                Title = "Mr",
                FirstName = "Karim",
                LastName = "Hossain",
                Gender = "Male",
                PhoneNumber = "+8801912345678",
                EmailAddress = "karim.hossain@email.com",
                Address = "Sylhet, Bangladesh",
                LeadSource = LeadSource.PhoneInquiry,
                LeadStatus = LeadStatus.Qualified,
                LeadPriority = LeadPriority.High,
                LeadScore = 85,
                DesiredInsuranceType = "Motor",
                DesiredCoverageAmount = 20000000, // 200,000 BDT
                QualificationStatus = QualificationStatus.Qualified,
                CreatedBy = "system"
            }
        };

        foreach (var lead in sampleLeads)
        {
            lead.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            lead.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            _leads[lead.LeadId] = lead;
        }

        _logger.LogInformation("Seeded {Count} sample leads", _leads.Count);
    }
}
