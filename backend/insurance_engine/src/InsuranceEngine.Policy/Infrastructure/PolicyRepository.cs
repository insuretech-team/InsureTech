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

public class PolicyRepository : IPolicyRepository
{
    private readonly PolicyDbContext _context;

    public PolicyRepository(PolicyDbContext context)
    {
        _context = context;
    }

    public async Task<PolicyEntity?> GetByIdAsync(Guid id) =>
        await _context.Policies
            .Include(p => p.Riders)
            .FirstOrDefaultAsync(p => p.Id == id);

    public async Task<PolicyEntity?> GetByIdWithNomineesAsync(Guid id) =>
        await _context.Policies
            .Include(p => p.Nominees)
            .Include(p => p.Riders)
            .FirstOrDefaultAsync(p => p.Id == id);

    public async Task<(List<PolicyEntity> Items, int TotalCount)> ListAsync(
        Guid? customerId, PolicyStatus? status, Guid? productId, int page, int pageSize)
    {
        var query = _context.Policies.AsQueryable();

        if (customerId.HasValue)
            query = query.Where(p => p.CustomerId == customerId.Value);
        if (status.HasValue)
            query = query.Where(p => p.Status == status.Value);
        if (productId.HasValue)
            query = query.Where(p => p.ProductId == productId.Value);

        var totalCount = await query.CountAsync();

        var items = await query
            .OrderByDescending(p => p.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync();

        return (items, totalCount);
    }

    public async Task<Guid> AddAsync(PolicyEntity policy)
    {
        _context.Policies.Add(policy);
        await _context.SaveChangesAsync();
        return policy.Id;
    }

    public async Task UpdateAsync(PolicyEntity policy)
    {
        _context.Policies.Update(policy);
        await _context.SaveChangesAsync();
    }

    public async Task<long> GetNextSequenceNumberAsync()
    {
        // Use PostgreSQL nextval to get thread-safe sequential number
        var connection = _context.Database.GetDbConnection();
        await using var command = connection.CreateCommand();
        command.CommandText = "SELECT nextval('insurance_schema.policy_number_seq')";

        if (connection.State != System.Data.ConnectionState.Open)
            await connection.OpenAsync();

        var result = await command.ExecuteScalarAsync();
        return Convert.ToInt64(result);
    }

    public async Task<string?> GetProductCodeAsync(Guid productId)
    {
        // Cross-schema query to products table
        var connection = _context.Database.GetDbConnection();
        await using var command = connection.CreateCommand();
        command.CommandText = "SELECT product_code FROM insurance_schema.products WHERE id = @productId LIMIT 1";
        var param = command.CreateParameter();
        param.ParameterName = "productId";
        param.Value = productId;
        command.Parameters.Add(param);

        if (connection.State != System.Data.ConnectionState.Open)
            await connection.OpenAsync();

        var result = await command.ExecuteScalarAsync();
        return result?.ToString();
    }
}
