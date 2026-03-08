using Google.Protobuf.WellKnownTypes;
using Insuretech.Policy.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface IRenewalReminderRepository
{
    Task<RenewalReminder> CreateAsync(RenewalReminder entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<RenewalReminder>> GetByStatusAsync(ReminderStatus status, CancellationToken cancellationToken = default);    Task<List<RenewalReminder>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<RenewalReminder> UpdateAsync(RenewalReminder entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class RenewalReminderRepository : IRenewalReminderRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<RenewalReminderRepository> _logger;

    public RenewalReminderRepository(PoliSyncDbContext context, ILogger<RenewalReminderRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<RenewalReminder> CreateAsync(RenewalReminder entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.ReminderId))
        {
            entity.ReminderId = Guid.NewGuid().ToString();
        }

        _context.RenewalReminders.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created RenewalReminder {Id}", entity.ReminderId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.RenewalReminders
            .FirstOrDefaultAsync(e => e.ReminderId == id, cancellationToken);
    }

    public async Task<List<RenewalReminder>> GetByStatusAsync(ReminderStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.RenewalReminders
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<RenewalReminder>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.RenewalReminders
            .OrderByDescending(e => e.ReminderId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<RenewalReminder> UpdateAsync(RenewalReminder entity, CancellationToken cancellationToken = default)
    {
        _context.RenewalReminders.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated RenewalReminder {Id}", entity.ReminderId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.RenewalReminders.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted RenewalReminder {Id}", id);
        }
    }
}
