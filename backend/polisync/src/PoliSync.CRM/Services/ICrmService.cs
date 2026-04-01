using Insuretech.Crm.Entity.V1;

namespace PoliSync.CRM.Services;

public interface ICrmService
{
    // Lead operations
    Task<Lead> CreateLeadAsync(Lead lead, CancellationToken cancellationToken = default);
    Task<Lead?> GetLeadAsync(string leadId, CancellationToken cancellationToken = default);
    Task<IEnumerable<Lead>> ListLeadsAsync(LeadStatus? status, LeadSource? source, string? assignedAgentId, CancellationToken cancellationToken = default);
    Task<Lead?> UpdateLeadAsync(string leadId, Lead lead, CancellationToken cancellationToken = default);
    Task<bool> DeleteLeadAsync(string leadId, bool permanent = false, CancellationToken cancellationToken = default);
    Task<Contact> ConvertLeadToContactAsync(string leadId, string conversionReason, CancellationToken cancellationToken = default);
    Task<Lead?> AssignLeadAsync(string leadId, string agentId, CancellationToken cancellationToken = default);

    // Contact operations
    Task<Contact> CreateContactAsync(Contact contact, CancellationToken cancellationToken = default);
    Task<Contact?> GetContactAsync(string contactId, CancellationToken cancellationToken = default);
    Task<IEnumerable<Contact>> ListContactsAsync(ContactStatus? status, ContactType? contactType, string? assignedAgentId, CancellationToken cancellationToken = default);
    Task<Contact?> UpdateContactAsync(string contactId, Contact contact, CancellationToken cancellationToken = default);
    Task<bool> DeleteContactAsync(string contactId, bool permanent = false, CancellationToken cancellationToken = default);

    // Statistics
    Task<(int TotalLeads, int NewLeads, int ContactedLeads, int QualifiedLeads, int ConvertedLeads, int LostLeads,
          int TotalContacts, int ActiveContacts, double ConversionRate, Dictionary<string, int> LeadsBySource)> 
        GetCrmStatisticsAsync(string? agentId, CancellationToken cancellationToken = default);

    Task<(int MyLeadsCount, int MyContactsCount, int NewLeadsToday, int LeadsConvertedThisMonth,
          List<Lead> RecentLeads, List<Contact> RecentContacts)> 
        GetAgentDashboardAsync(string agentId, CancellationToken cancellationToken = default);
}

public interface ILeadRepository
{
    Task<Lead?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<IEnumerable<Lead>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<Lead>> GetByFilterAsync(LeadStatus? status, LeadSource? source, string? assignedAgentId, CancellationToken cancellationToken = default);
    Task<IEnumerable<Lead>> GetByAssignedAgentAsync(string agentId, CancellationToken cancellationToken = default);
    Task<Lead> CreateAsync(Lead lead, CancellationToken cancellationToken = default);
    Task<Lead?> UpdateAsync(Lead lead, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default);
}

public interface IContactRepository
{
    Task<Contact?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<IEnumerable<Contact>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<Contact>> GetByFilterAsync(ContactStatus? status, ContactType? contactType, string? assignedAgentId, CancellationToken cancellationToken = default);
    Task<IEnumerable<Contact>> GetByAssignedAgentAsync(string agentId, CancellationToken cancellationToken = default);
    Task<Contact> CreateAsync(Contact contact, CancellationToken cancellationToken = default);
    Task<Contact?> UpdateAsync(Contact contact, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default);
}
