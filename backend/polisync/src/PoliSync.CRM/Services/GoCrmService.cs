using Insuretech.Crm.Entity.V1;
using Insuretech.Crm.Services.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.CRM.Services;

public sealed class GoCrmService : ICrmService
{
    private readonly InsuranceServiceClient _client;

    public GoCrmService(InsuranceServiceClient client) => _client = client;

    public async Task<Lead> CreateLeadAsync(Lead lead, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.CreateLeadAsync(
            new CreateLeadRequest
            {
                Title = lead.Title,
                FirstName = lead.FirstName,
                LastName = lead.LastName,
                DateOfBirth = lead.DateOfBirth,
                Gender = lead.Gender,
                PhoneNumber = lead.PhoneNumber,
                EmailAddress = lead.EmailAddress,
                Address = lead.Address,
                LeadSource = lead.LeadSource,
                LeadPriority = lead.LeadPriority,
                AssignedAgentId = lead.AssignedAgentId,
                DesiredInsuranceType = lead.DesiredInsuranceType,
                DesiredCoverageAmount = lead.DesiredCoverageAmount,
                SpecificRequirements = lead.SpecificRequirements,
                Metadata = { lead.Metadata },
                CreatedBy = lead.CreatedBy
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Lead;
    }

    public async Task<Lead?> GetLeadAsync(string leadId, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.GetLeadAsync(
            new GetLeadRequest { LeadId = leadId },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Lead : null;
    }

    public async Task<IEnumerable<Lead>> ListLeadsAsync(LeadStatus? status, LeadSource? source, string? assignedAgentId, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.ListLeadsAsync(
            new ListLeadsRequest
            {
                Status = status ?? LeadStatus.Unspecified,
                Source = source ?? LeadSource.Unspecified,
                AssignedAgentId = assignedAgentId ?? string.Empty,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Leads;
    }

    public async Task<Lead?> UpdateLeadAsync(string leadId, Lead lead, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.UpdateLeadAsync(
            new UpdateLeadRequest
            {
                LeadId = leadId,
                Title = lead.Title,
                FirstName = lead.FirstName,
                LastName = lead.LastName,
                PhoneNumber = lead.PhoneNumber,
                EmailAddress = lead.EmailAddress,
                Address = lead.Address,
                LeadStatus = lead.LeadStatus,
                LeadPriority = lead.LeadPriority,
                LeadScore = lead.LeadScore,
                AssignedAgentId = lead.AssignedAgentId,
                DesiredInsuranceType = lead.DesiredInsuranceType,
                DesiredCoverageAmount = lead.DesiredCoverageAmount,
                QualificationStatus = lead.QualificationStatus,
                Metadata = { lead.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Lead : null;
    }

    public async Task<bool> DeleteLeadAsync(string leadId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.DeleteLeadAsync(
            new DeleteLeadRequest { LeadId = leadId, Permanent = permanent },
            _client.BuildCallOptions(cancellationToken));
        return response.Success && response.Error is null;
    }

    public async Task<Contact> ConvertLeadToContactAsync(string leadId, string conversionReason, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.ConvertLeadToContactAsync(
            new ConvertLeadToContactRequest
            {
                LeadId = leadId,
                ConversionReason = conversionReason
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Contact;
    }

    public async Task<Lead?> AssignLeadAsync(string leadId, string agentId, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.AssignLeadAsync(
            new AssignLeadRequest
            {
                LeadId = leadId,
                AgentId = agentId
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Lead : null;
    }

    public async Task<Contact> CreateContactAsync(Contact contact, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.CreateContactAsync(
            new CreateContactRequest
            {
                FirstName = contact.FirstName,
                LastName = contact.LastName,
                DateOfBirth = contact.DateOfBirth,
                Gender = contact.Gender,
                PhoneNumber = contact.PhoneNumber,
                EmailAddress = contact.EmailAddress,
                Address = contact.Address,
                PreferredContactMethod = contact.PreferredContactMethod,
                ReferralSource = contact.ReferralSource,
                PreferredLanguage = contact.PreferredLanguage,
                MarketingConsent = contact.MarketingConsent,
                AssignedAgentId = contact.AssignedAgentId,
                ContactType = contact.ContactType,
                Metadata = { contact.Metadata },
                CreatedBy = contact.CreatedBy
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Contact;
    }

    public async Task<Contact?> GetContactAsync(string contactId, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.GetContactAsync(
            new GetContactRequest { ContactId = contactId },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Contact : null;
    }

    public async Task<IEnumerable<Contact>> ListContactsAsync(ContactStatus? status, ContactType? contactType, string? assignedAgentId, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.ListContactsAsync(
            new ListContactsRequest
            {
                Status = status ?? ContactStatus.Unspecified,
                ContactType = contactType ?? ContactType.Unspecified,
                AssignedAgentId = assignedAgentId ?? string.Empty,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Contacts;
    }

    public async Task<Contact?> UpdateContactAsync(string contactId, Contact contact, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.UpdateContactAsync(
            new UpdateContactRequest
            {
                ContactId = contactId,
                FirstName = contact.FirstName,
                LastName = contact.LastName,
                PhoneNumber = contact.PhoneNumber,
                EmailAddress = contact.EmailAddress,
                Address = contact.Address,
                PreferredContactMethod = contact.PreferredContactMethod,
                MarketingConsent = contact.MarketingConsent,
                AssignedAgentId = contact.AssignedAgentId,
                ContactStatus = contact.ContactStatus,
                Metadata = { contact.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Contact : null;
    }

    public async Task<bool> DeleteContactAsync(string contactId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.DeleteContactAsync(
            new DeleteContactRequest { ContactId = contactId, Permanent = permanent },
            _client.BuildCallOptions(cancellationToken));
        return response.Success && response.Error is null;
    }

    public async Task<(int TotalLeads, int NewLeads, int ContactedLeads, int QualifiedLeads, int ConvertedLeads, int LostLeads, int TotalContacts, int ActiveContacts, double ConversionRate, Dictionary<string, int> LeadsBySource)> GetCrmStatisticsAsync(string? agentId, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.GetCrmStatisticsAsync(
            new GetCrmStatisticsRequest { AgentId = agentId ?? string.Empty },
            _client.BuildCallOptions(cancellationToken));

        return (
            response.TotalLeads,
            response.NewLeads,
            response.ContactedLeads,
            response.QualifiedLeads,
            response.ConvertedLeads,
            response.LostLeads,
            response.TotalContacts,
            response.ActiveContacts,
            response.ConversionRate,
            response.LeadsBySource.ToDictionary(kvp => kvp.Key, kvp => kvp.Value));
    }

    public async Task<(int MyLeadsCount, int MyContactsCount, int NewLeadsToday, int LeadsConvertedThisMonth, List<Lead> RecentLeads, List<Contact> RecentContacts)> GetAgentDashboardAsync(string agentId, CancellationToken cancellationToken = default)
    {
        var response = await _client.CrmClient.GetAgentDashboardAsync(
            new GetAgentDashboardRequest { AgentId = agentId },
            _client.BuildCallOptions(cancellationToken));

        return (
            response.MyLeadsCount,
            response.MyContactsCount,
            response.NewLeadsToday,
            response.LeadsConvertedThisMonth,
            response.RecentLeads.ToList(),
            response.RecentContacts.ToList());
    }
}
