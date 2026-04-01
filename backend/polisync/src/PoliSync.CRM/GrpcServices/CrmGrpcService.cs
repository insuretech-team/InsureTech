using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Crm.Entity.V1;
using Insuretech.Crm.Services.V1;
using PoliSync.CRM.Services;
using ProtoCrmService = Insuretech.Crm.Services.V1.CrmService;

namespace PoliSync.CRM.GrpcServices;

public class CrmGrpcService : ProtoCrmService.CrmServiceBase
{
    private readonly ICrmService _crmService;
    private readonly ILogger<CrmGrpcService> _logger;

    public CrmGrpcService(
        ICrmService crmService,
        ILogger<CrmGrpcService> logger)
    {
        _crmService = crmService;
        _logger = logger;
    }

    // Lead operations
    public override async Task<CreateLeadResponse> CreateLead(
        CreateLeadRequest request,
        ServerCallContext context)
    {
        try
        {
            var lead = new Lead
            {
                Title = request.Title,
                FirstName = request.FirstName,
                LastName = request.LastName,
                DateOfBirth = request.DateOfBirth,
                Gender = request.Gender,
                PhoneNumber = request.PhoneNumber,
                EmailAddress = request.EmailAddress,
                Address = request.Address,
                LeadSource = request.LeadSource,
                LeadPriority = request.LeadPriority,
                AssignedAgentId = request.AssignedAgentId,
                DesiredInsuranceType = request.DesiredInsuranceType,
                DesiredCoverageAmount = request.DesiredCoverageAmount,
                SpecificRequirements = request.SpecificRequirements,
                Metadata = { request.Metadata },
                CreatedBy = request.CreatedBy
            };

            var created = await _crmService.CreateLeadAsync(lead, context.CancellationToken);

            return new CreateLeadResponse { Lead = created };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating lead");
            return new CreateLeadResponse
            {
                Error = new Error
                {
                    Code = "CREATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<GetLeadResponse> GetLead(
        GetLeadRequest request,
        ServerCallContext context)
    {
        var lead = await _crmService.GetLeadAsync(request.LeadId, context.CancellationToken);
        
        if (lead == null)
        {
            return new GetLeadResponse
            {
                Error = new Error
                {
                    Code = "LEAD_NOT_FOUND",
                    Message = $"Lead {request.LeadId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetLeadResponse { Lead = lead };
    }

    public override async Task<ListLeadsResponse> ListLeads(
        ListLeadsRequest request,
        ServerCallContext context)
    {
        var leads = await _crmService.ListLeadsAsync(
            request.Status != LeadStatus.Unspecified ? request.Status : null,
            request.Source != LeadSource.Unspecified ? request.Source : null,
            string.IsNullOrEmpty(request.AssignedAgentId) ? null : request.AssignedAgentId,
            context.CancellationToken);

        var leadList = leads.ToList();

        return new ListLeadsResponse
        {
            Leads = { leadList },
            TotalCount = leadList.Count
        };
    }

    public override async Task<UpdateLeadResponse> UpdateLead(
        UpdateLeadRequest request,
        ServerCallContext context)
    {
        try
        {
            var lead = new Lead
            {
                LeadId = request.LeadId,
                Title = request.Title,
                FirstName = request.FirstName,
                LastName = request.LastName,
                PhoneNumber = request.PhoneNumber,
                EmailAddress = request.EmailAddress,
                Address = request.Address,
                LeadStatus = request.LeadStatus,
                LeadPriority = request.LeadPriority,
                LeadScore = request.LeadScore,
                AssignedAgentId = request.AssignedAgentId,
                DesiredInsuranceType = request.DesiredInsuranceType,
                DesiredCoverageAmount = request.DesiredCoverageAmount,
                QualificationStatus = request.QualificationStatus,
                Metadata = { request.Metadata }
            };

            var updated = await _crmService.UpdateLeadAsync(request.LeadId, lead, context.CancellationToken);

            if (updated == null)
            {
                return new UpdateLeadResponse
                {
                    Error = new Error
                    {
                        Code = "LEAD_NOT_FOUND",
                        Message = $"Lead {request.LeadId} not found",
                        HttpStatusCode = 404
                    }
                };
            }

            return new UpdateLeadResponse { Lead = updated };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating lead");
            return new UpdateLeadResponse
            {
                Error = new Error
                {
                    Code = "UPDATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<DeleteLeadResponse> DeleteLead(
        DeleteLeadRequest request,
        ServerCallContext context)
    {
        var result = await _crmService.DeleteLeadAsync(
            request.LeadId,
            request.Permanent,
            context.CancellationToken);

        return new DeleteLeadResponse { Success = result };
    }

    public override async Task<ConvertLeadToContactResponse> ConvertLeadToContact(
        ConvertLeadToContactRequest request,
        ServerCallContext context)
    {
        try
        {
            var contact = await _crmService.ConvertLeadToContactAsync(
                request.LeadId,
                request.ConversionReason,
                context.CancellationToken);

            var lead = await _crmService.GetLeadAsync(request.LeadId, context.CancellationToken);

            return new ConvertLeadToContactResponse
            {
                Contact = contact,
                Lead = lead
            };
        }
        catch (InvalidOperationException ex)
        {
            return new ConvertLeadToContactResponse
            {
                Error = new Error
                {
                    Code = "CONVERSION_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 400
                }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error converting lead to contact");
            return new ConvertLeadToContactResponse
            {
                Error = new Error
                {
                    Code = "CONVERSION_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<AssignLeadResponse> AssignLead(
        AssignLeadRequest request,
        ServerCallContext context)
    {
        var lead = await _crmService.AssignLeadAsync(
            request.LeadId,
            request.AgentId,
            context.CancellationToken);

        if (lead == null)
        {
            return new AssignLeadResponse
            {
                Error = new Error
                {
                    Code = "LEAD_NOT_FOUND",
                    Message = $"Lead {request.LeadId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new AssignLeadResponse { Lead = lead };
    }

    // Contact operations
    public override async Task<CreateContactResponse> CreateContact(
        CreateContactRequest request,
        ServerCallContext context)
    {
        try
        {
            var contact = new Contact
            {
                FirstName = request.FirstName,
                LastName = request.LastName,
                DateOfBirth = request.DateOfBirth,
                Gender = request.Gender,
                PhoneNumber = request.PhoneNumber,
                EmailAddress = request.EmailAddress,
                Address = request.Address,
                PreferredContactMethod = request.PreferredContactMethod,
                ReferralSource = request.ReferralSource,
                PreferredLanguage = request.PreferredLanguage,
                MarketingConsent = request.MarketingConsent,
                AssignedAgentId = request.AssignedAgentId,
                ContactType = request.ContactType,
                Metadata = { request.Metadata },
                CreatedBy = request.CreatedBy
            };

            var created = await _crmService.CreateContactAsync(contact, context.CancellationToken);

            return new CreateContactResponse { Contact = created };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating contact");
            return new CreateContactResponse
            {
                Error = new Error
                {
                    Code = "CREATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<GetContactResponse> GetContact(
        GetContactRequest request,
        ServerCallContext context)
    {
        var contact = await _crmService.GetContactAsync(request.ContactId, context.CancellationToken);
        
        if (contact == null)
        {
            return new GetContactResponse
            {
                Error = new Error
                {
                    Code = "CONTACT_NOT_FOUND",
                    Message = $"Contact {request.ContactId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetContactResponse { Contact = contact };
    }

    public override async Task<ListContactsResponse> ListContacts(
        ListContactsRequest request,
        ServerCallContext context)
    {
        var contacts = await _crmService.ListContactsAsync(
            request.Status != ContactStatus.Unspecified ? request.Status : null,
            request.ContactType != ContactType.Unspecified ? request.ContactType : null,
            string.IsNullOrEmpty(request.AssignedAgentId) ? null : request.AssignedAgentId,
            context.CancellationToken);

        var contactList = contacts.ToList();

        return new ListContactsResponse
        {
            Contacts = { contactList },
            TotalCount = contactList.Count
        };
    }

    public override async Task<UpdateContactResponse> UpdateContact(
        UpdateContactRequest request,
        ServerCallContext context)
    {
        try
        {
            var contact = new Contact
            {
                ContactId = request.ContactId,
                FirstName = request.FirstName,
                LastName = request.LastName,
                PhoneNumber = request.PhoneNumber,
                EmailAddress = request.EmailAddress,
                Address = request.Address,
                PreferredContactMethod = request.PreferredContactMethod,
                MarketingConsent = request.MarketingConsent,
                AssignedAgentId = request.AssignedAgentId,
                ContactStatus = request.ContactStatus,
                Metadata = { request.Metadata }
            };

            var updated = await _crmService.UpdateContactAsync(request.ContactId, contact, context.CancellationToken);

            if (updated == null)
            {
                return new UpdateContactResponse
                {
                    Error = new Error
                    {
                        Code = "CONTACT_NOT_FOUND",
                        Message = $"Contact {request.ContactId} not found",
                        HttpStatusCode = 404
                    }
                };
            }

            return new UpdateContactResponse { Contact = updated };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating contact");
            return new UpdateContactResponse
            {
                Error = new Error
                {
                    Code = "UPDATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<DeleteContactResponse> DeleteContact(
        DeleteContactRequest request,
        ServerCallContext context)
    {
        var result = await _crmService.DeleteContactAsync(
            request.ContactId,
            request.Permanent,
            context.CancellationToken);

        return new DeleteContactResponse { Success = result };
    }

    // Statistics
    public override async Task<GetCrmStatisticsResponse> GetCrmStatistics(
        GetCrmStatisticsRequest request,
        ServerCallContext context)
    {
        var (totalLeads, newLeads, contactedLeads, qualifiedLeads, convertedLeads, lostLeads,
             totalContacts, activeContacts, conversionRate, leadsBySource) = 
            await _crmService.GetCrmStatisticsAsync(
                string.IsNullOrEmpty(request.AgentId) ? null : request.AgentId,
                context.CancellationToken);

        return new GetCrmStatisticsResponse
        {
            TotalLeads = totalLeads,
            NewLeads = newLeads,
            ContactedLeads = contactedLeads,
            QualifiedLeads = qualifiedLeads,
            ConvertedLeads = convertedLeads,
            LostLeads = lostLeads,
            TotalContacts = totalContacts,
            ActiveContacts = activeContacts,
            ConversionRate = conversionRate,
            LeadsBySource = { leadsBySource }
        };
    }

    public override async Task<GetAgentDashboardResponse> GetAgentDashboard(
        GetAgentDashboardRequest request,
        ServerCallContext context)
    {
        var (myLeadsCount, myContactsCount, newLeadsToday, leadsConvertedThisMonth,
             recentLeads, recentContacts) = 
            await _crmService.GetAgentDashboardAsync(request.AgentId, context.CancellationToken);

        return new GetAgentDashboardResponse
        {
            MyLeadsCount = myLeadsCount,
            MyContactsCount = myContactsCount,
            NewLeadsToday = newLeadsToday,
            LeadsConvertedThisMonth = leadsConvertedThisMonth,
            RecentLeads = { recentLeads },
            RecentContacts = { recentContacts }
        };
    }
}
