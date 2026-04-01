using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Crm.Entity.V1;

namespace PoliSync.CRM.Services;

public class CrmService : ICrmService
{
    private readonly ILeadRepository _leadRepository;
    private readonly IContactRepository _contactRepository;
    private readonly ILogger<CrmService> _logger;

    public CrmService(
        ILeadRepository leadRepository,
        IContactRepository contactRepository,
        ILogger<CrmService> logger)
    {
        _leadRepository = leadRepository;
        _contactRepository = contactRepository;
        _logger = logger;
    }

    // Lead operations
    public Task<Lead> CreateLeadAsync(Lead lead, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Creating lead: {FirstName} {LastName}", lead.FirstName, lead.LastName);
        return _leadRepository.CreateAsync(lead, cancellationToken);
    }

    public Task<Lead?> GetLeadAsync(string leadId, CancellationToken cancellationToken = default)
    {
        return _leadRepository.GetByIdAsync(leadId, cancellationToken);
    }

    public Task<IEnumerable<Lead>> ListLeadsAsync(LeadStatus? status, LeadSource? source, string? assignedAgentId, CancellationToken cancellationToken = default)
    {
        return _leadRepository.GetByFilterAsync(status, source, assignedAgentId, cancellationToken);
    }

    public Task<Lead?> UpdateLeadAsync(string leadId, Lead lead, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Updating lead: {LeadId}", leadId);
        lead.LeadId = leadId;
        return _leadRepository.UpdateAsync(lead, cancellationToken);
    }

    public Task<bool> DeleteLeadAsync(string leadId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Deleting lead: {LeadId} (permanent: {Permanent})", leadId, permanent);
        return _leadRepository.DeleteAsync(leadId, permanent, cancellationToken);
    }

    public async Task<Contact> ConvertLeadToContactAsync(string leadId, string conversionReason, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Converting lead {LeadId} to contact", leadId);

        var lead = await _leadRepository.GetByIdAsync(leadId, cancellationToken);
        if (lead == null)
        {
            throw new InvalidOperationException($"Lead {leadId} not found");
        }

        if (lead.LeadStatus != LeadStatus.New && 
            lead.LeadStatus != LeadStatus.Contacted && 
            lead.LeadStatus != LeadStatus.Qualified)
        {
            throw new InvalidOperationException($"Lead {leadId} cannot be converted. Current status: {lead.LeadStatus}");
        }

        // Create contact from lead
        var contact = new Contact
        {
            ContactId = Guid.NewGuid().ToString(),
            FirstName = lead.FirstName,
            LastName = lead.LastName,
            DateOfBirth = lead.DateOfBirth,
            Gender = lead.Gender,
            PhoneNumber = lead.PhoneNumber,
            EmailAddress = lead.EmailAddress,
            Address = lead.Address,
            AlternatePhoneNumber = lead.AlternatePhoneNumber,
            AdditionalEmailAddress = lead.AdditionalEmailAddress,
            AssignedAgentId = lead.AssignedAgentId,
            ContactType = ContactType.Individual,
            ContactStatus = ContactStatus.Active,
            ReferralSource = lead.LeadSource.ToString(),
            PreferredLanguage = "English",
            CreatedBy = lead.CreatedBy
        };

        var createdContact = await _contactRepository.CreateAsync(contact, cancellationToken);

        // Update lead as converted
        lead.LeadStatus = LeadStatus.Converted;
        lead.ConvertedContactId = createdContact.ContactId;
        lead.ConvertedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        lead.ConversionReason = conversionReason;
        await _leadRepository.UpdateAsync(lead, cancellationToken);

        _logger.LogInformation("Lead {LeadId} converted to contact {ContactId}", leadId, createdContact.ContactId);

        return createdContact;
    }

    public async Task<Lead?> AssignLeadAsync(string leadId, string agentId, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Assigning lead {LeadId} to agent {AgentId}", leadId, agentId);

        var lead = await _leadRepository.GetByIdAsync(leadId, cancellationToken);
        if (lead == null)
        {
            return null;
        }

        lead.AssignedAgentId = agentId;
        return await _leadRepository.UpdateAsync(lead, cancellationToken);
    }

    // Contact operations
    public Task<Contact> CreateContactAsync(Contact contact, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Creating contact: {FirstName} {LastName}", contact.FirstName, contact.LastName);
        return _contactRepository.CreateAsync(contact, cancellationToken);
    }

    public Task<Contact?> GetContactAsync(string contactId, CancellationToken cancellationToken = default)
    {
        return _contactRepository.GetByIdAsync(contactId, cancellationToken);
    }

    public Task<IEnumerable<Contact>> ListContactsAsync(ContactStatus? status, ContactType? contactType, string? assignedAgentId, CancellationToken cancellationToken = default)
    {
        return _contactRepository.GetByFilterAsync(status, contactType, assignedAgentId, cancellationToken);
    }

    public Task<Contact?> UpdateContactAsync(string contactId, Contact contact, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Updating contact: {ContactId}", contactId);
        contact.ContactId = contactId;
        return _contactRepository.UpdateAsync(contact, cancellationToken);
    }

    public Task<bool> DeleteContactAsync(string contactId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Deleting contact: {ContactId} (permanent: {Permanent})", contactId, permanent);
        return _contactRepository.DeleteAsync(contactId, permanent, cancellationToken);
    }

    // Statistics
    public async Task<(int TotalLeads, int NewLeads, int ContactedLeads, int QualifiedLeads, int ConvertedLeads, int LostLeads,
          int TotalContacts, int ActiveContacts, double ConversionRate, Dictionary<string, int> LeadsBySource)> 
        GetCrmStatisticsAsync(string? agentId, CancellationToken cancellationToken = default)
    {
        var leads = await _leadRepository.GetAllAsync(cancellationToken);
        var contacts = await _contactRepository.GetAllAsync(cancellationToken);

        // Filter by agent if specified
        if (!string.IsNullOrEmpty(agentId))
        {
            leads = leads.Where(l => l.AssignedAgentId == agentId);
            contacts = contacts.Where(c => c.AssignedAgentId == agentId);
        }

        var leadList = leads.ToList();
        var totalLeads = leadList.Count;
        var newLeads = leadList.Count(l => l.LeadStatus == LeadStatus.New);
        var contactedLeads = leadList.Count(l => l.LeadStatus == LeadStatus.Contacted);
        var qualifiedLeads = leadList.Count(l => l.LeadStatus == LeadStatus.Qualified);
        var convertedLeads = leadList.Count(l => l.LeadStatus == LeadStatus.Converted);
        var lostLeads = leadList.Count(l => l.LeadStatus == LeadStatus.Lost);

        var contactList = contacts.ToList();
        var totalContacts = contactList.Count;
        var activeContacts = contactList.Count(c => c.ContactStatus == ContactStatus.Active);

        var conversionRate = totalLeads > 0 ? (double)convertedLeads / totalLeads * 100 : 0;

        var leadsBySource = leadList
            .GroupBy(l => l.LeadSource.ToString())
            .ToDictionary(g => g.Key, g => g.Count());

        return (totalLeads, newLeads, contactedLeads, qualifiedLeads, convertedLeads, lostLeads,
                totalContacts, activeContacts, conversionRate, leadsBySource);
    }

    public async Task<(int MyLeadsCount, int MyContactsCount, int NewLeadsToday, int LeadsConvertedThisMonth,
          List<Lead> RecentLeads, List<Contact> RecentContacts)> 
        GetAgentDashboardAsync(string agentId, CancellationToken cancellationToken = default)
    {
        var myLeads = await _leadRepository.GetByAssignedAgentAsync(agentId, cancellationToken);
        var myContacts = await _contactRepository.GetByAssignedAgentAsync(agentId, cancellationToken);

        var today = DateTime.UtcNow.Date;
        var startOfMonth = new DateTime(today.Year, today.Month, 1);

        var myLeadsList = myLeads.ToList();
        var myContactsList = myContacts.ToList();

        var newLeadsToday = myLeadsList.Count(l => l.CreatedAt.ToDateTime().Date == today);
        var leadsConvertedThisMonth = myLeadsList.Count(l => 
            l.LeadStatus == LeadStatus.Converted && 
            l.ConvertedAt != null && 
            l.ConvertedAt.ToDateTime() >= startOfMonth);

        var recentLeads = myLeadsList
            .OrderByDescending(l => l.CreatedAt)
            .Take(5)
            .ToList();

        var recentContacts = myContactsList
            .OrderByDescending(c => c.CreatedAt)
            .Take(5)
            .ToList();

        return (myLeadsList.Count, myContactsList.Count, newLeadsToday, leadsConvertedThisMonth,
                recentLeads, recentContacts);
    }
}
