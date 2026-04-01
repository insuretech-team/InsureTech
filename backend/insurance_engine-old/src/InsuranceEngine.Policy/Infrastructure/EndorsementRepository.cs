using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Infrastructure.Persistence;

namespace InsuranceEngine.Policy.Infrastructure;

public class EndorsementRepository : IEndorsementRepository
{
    private readonly PolicyDbContext _context;

    public EndorsementRepository(PolicyDbContext context)
    {
        _context = context;
    }

    public async Task<Endorsement?> GetByIdAsync(Guid id) =>
        await _context.Endorsements.FirstOrDefaultAsync(e => e.Id == id);

    public async Task<(List<Endorsement> Items, int TotalCount)> ListAsync(
        Guid? policyId, EndorsementStatus? status, int page, int pageSize)
    {
        var query = _context.Endorsements.AsQueryable();

        if (policyId.HasValue)
            query = query.Where(e => e.PolicyId == policyId.Value);
        if (status.HasValue)
            query = query.Where(e => e.Status == status.Value);

        var totalCount = await query.CountAsync();

        var items = await query
            .OrderByDescending(e => e.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync();

        return (items, totalCount);
    }

    public async Task<Guid> AddAsync(Endorsement endorsement)
    {
        _context.Endorsements.Add(endorsement);
        await _context.SaveChangesAsync();
        return endorsement.Id;
    }

    public async Task UpdateAsync(Endorsement endorsement)
    {
        _context.Endorsements.Update(endorsement);
        await _context.SaveChangesAsync();
    }

    public async Task<long> GetNextSequenceNumberAsync()
    {
        var connection = _context.Database.GetDbConnection();
        await using var command = connection.CreateCommand();
        command.CommandText = "SELECT nextval('insurance_schema.endorsement_number_seq')";

        if (connection.State != System.Data.ConnectionState.Open)
            await connection.OpenAsync();

        var result = await command.ExecuteScalarAsync();
        return Convert.ToInt64(result);
    }
}
