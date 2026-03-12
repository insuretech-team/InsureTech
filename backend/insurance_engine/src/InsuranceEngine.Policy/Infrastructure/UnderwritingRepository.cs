using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Policy.Infrastructure;

public class UnderwritingRepository : IUnderwritingRepository
{
    private readonly PolicyDbContext _dbContext;

    public UnderwritingRepository(PolicyDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    public async Task<Quote?> GetQuoteByIdAsync(Guid id)
    {
        return await _dbContext.Quotes
            .FirstOrDefaultAsync(q => q.Id == id);
    }

    public async Task<Quote?> GetQuoteByNumberAsync(string quoteNumber)
    {
        return await _dbContext.Quotes
            .FirstOrDefaultAsync(q => q.QuoteNumber == quoteNumber);
    }

    public async Task<(List<Quote> Items, int TotalCount)> ListQuotesAsync(
        Guid? beneficiaryId, QuoteStatus? status, int page, int pageSize)
    {
        var query = _dbContext.Quotes.AsQueryable();

        if (beneficiaryId.HasValue)
            query = query.Where(q => q.BeneficiaryId == beneficiaryId.Value);

        if (status.HasValue)
            query = query.Where(q => q.Status == status.Value);

        var totalCount = await query.CountAsync();
        var items = await query
            .OrderByDescending(q => q.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync();

        return (items, totalCount);
    }

    public async Task<Guid> AddQuoteAsync(Quote quote)
    {
        _dbContext.Quotes.Add(quote);
        await _dbContext.SaveChangesAsync();
        return quote.Id;
    }

    public async Task UpdateQuoteAsync(Quote quote)
    {
        _dbContext.Quotes.Update(quote);
        await _dbContext.SaveChangesAsync();
    }

    public async Task<UnderwritingHealthDeclaration?> GetHealthDeclarationByIdAsync(Guid id)
    {
        return await _dbContext.HealthDeclarations
            .FirstOrDefaultAsync(h => h.Id == id);
    }

    public async Task<UnderwritingHealthDeclaration?> GetHealthDeclarationByQuoteIdAsync(Guid quoteId)
    {
        return await _dbContext.HealthDeclarations
            .FirstOrDefaultAsync(h => h.QuoteId == quoteId);
    }

    public async Task AddHealthDeclarationAsync(UnderwritingHealthDeclaration declaration)
    {
        _dbContext.HealthDeclarations.Add(declaration);
        await _dbContext.SaveChangesAsync();
    }

    public async Task UpdateHealthDeclarationAsync(UnderwritingHealthDeclaration declaration)
    {
        _dbContext.HealthDeclarations.Update(declaration);
        await _dbContext.SaveChangesAsync();
    }

    public async Task<UnderwritingDecision?> GetDecisionByQuoteIdAsync(Guid quoteId)
    {
        return await _dbContext.UnderwritingDecisions
            .OrderByDescending(d => d.DecidedAt)
            .FirstOrDefaultAsync(d => d.QuoteId == quoteId);
    }

    public async Task AddDecisionAsync(UnderwritingDecision decision)
    {
        _dbContext.UnderwritingDecisions.Add(decision);
        await _dbContext.SaveChangesAsync();
    }

    public async Task UpdateDecisionAsync(UnderwritingDecision decision)
    {
        _dbContext.UnderwritingDecisions.Update(decision);
        await _dbContext.SaveChangesAsync();
    }

    public async Task<(List<UnderwritingDecision> Items, int TotalCount)> GetDecisionHistoryAsync(Guid quoteId)
    {
        var query = _dbContext.UnderwritingDecisions
            .Where(d => d.QuoteId == quoteId);

        var totalCount = await query.CountAsync();
        var items = await query
            .OrderByDescending(d => d.DecidedAt)
            .ToListAsync();

        return (items, totalCount);
    }

    public async Task<long> GetNextQuoteSequenceAsync()
    {
        var connection = _dbContext.Database.GetDbConnection();
        await connection.OpenAsync();
        using var command = connection.CreateCommand();
        command.CommandText = "SELECT nextval('insurance_schema.quote_number_seq')";
        var result = await command.ExecuteScalarAsync();
        return (long)(result ?? 0);
    }
}
