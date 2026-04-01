using System.Collections.Concurrent;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Crm.Entity.V1;

namespace PoliSync.CRM.Services;

public class InMemoryContactRepository : IContactRepository
{
    private readonly ConcurrentDictionary<string, Contact> _contacts = new();
    private readonly ILogger<InMemoryContactRepository> _logger;

    public InMemoryContactRepository(ILogger<InMemoryContactRepository> logger)
    {
        _logger = logger;
    }

    public Task<Contact?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _contacts.TryGetValue(id, out var contact);
        return Task.FromResult(contact);
    }

    public Task<IEnumerable<Contact>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_contacts.Values.AsEnumerable());
    }

    public Task<IEnumerable<Contact>> GetByFilterAsync(ContactStatus? status, ContactType? contactType, string? assignedAgentId, CancellationToken cancellationToken = default)
    {
        var query = _contacts.Values.AsEnumerable();
        
        if (status.HasValue)
        {
            query = query.Where(c => c.ContactStatus == status.Value);
        }
        
        if (contactType.HasValue)
        {
            query = query.Where(c => c.ContactType == contactType.Value);
        }
        
        if (!string.IsNullOrEmpty(assignedAgentId))
        {
            query = query.Where(c => c.AssignedAgentId == assignedAgentId);
        }
        
        return Task.FromResult<IEnumerable<Contact>>(query.OrderByDescending(c => c.CreatedAt));
    }

    public Task<IEnumerable<Contact>> GetByAssignedAgentAsync(string agentId, CancellationToken cancellationToken = default)
    {
        var contacts = _contacts.Values
            .Where(c => c.AssignedAgentId == agentId)
            .OrderByDescending(c => c.CreatedAt)
            .AsEnumerable();
        return Task.FromResult(contacts);
    }

    public Task<Contact> CreateAsync(Contact contact, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(contact.ContactId))
        {
            contact.ContactId = Guid.NewGuid().ToString();
        }
        
        contact.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        contact.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        _contacts[contact.ContactId] = contact;
        _logger.LogInformation("Created contact: {ContactId} - {FirstName} {LastName}", 
            contact.ContactId, contact.FirstName, contact.LastName);
        
        return Task.FromResult(contact);
    }

    public Task<Contact?> UpdateAsync(Contact contact, CancellationToken cancellationToken = default)
    {
        if (!_contacts.ContainsKey(contact.ContactId))
        {
            return Task.FromResult<Contact?>(null);
        }

        contact.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _contacts[contact.ContactId] = contact;
        
        _logger.LogInformation("Updated contact: {ContactId}", contact.ContactId);
        
        return Task.FromResult<Contact?>(contact);
    }

    public Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default)
    {
        if (permanent)
        {
            var result = _contacts.TryRemove(id, out _);
            if (result)
            {
                _logger.LogInformation("Permanently deleted contact: {ContactId}", id);
            }
            return Task.FromResult(result);
        }
        else
        {
            if (_contacts.TryGetValue(id, out var contact))
            {
                contact.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                contact.ContactStatus = ContactStatus.Archived;
                _logger.LogInformation("Soft deleted contact: {ContactId}", id);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }
}
