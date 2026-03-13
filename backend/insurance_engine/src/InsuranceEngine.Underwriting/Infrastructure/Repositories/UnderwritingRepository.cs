using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.Underwriting.Domain.Entities;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.Underwriting.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Underwriting.Infrastructure.Repositories;

public class UnderwritingRepository : IUnderwritingRepository
{
    private readonly UnderwritingDbContext _context;

    public UnderwritingRepository(UnderwritingDbContext context)
    {
        _context = context;
    }

    public async Task<Quote?> GetQuoteByIdAsync(Guid id)
    {
        return await _context.Quotes.FirstOrDefaultAsync(q => q.Id == id);
    }

    public async Task<Quote?> GetQuoteByNumberAsync(string quoteNumber)
    {
        return await _context.Quotes.FirstOrDefaultAsync(q => q.QuoteNumber == quoteNumber);
    }

    public async Task<(List<Quote> Items, int TotalCount)> ListQuotesAsync(
        Guid? beneficiaryId, QuoteStatus? status, int page, int pageSize)
    {
        var query = _context.Quotes.AsQueryable();
        if (beneficiaryId.HasValue) query = query.Where(q => q.BeneficiaryId == beneficiaryId);
        if (status.HasValue) query = query.Where(q => q.Status == status);

        var total = await query.CountAsync();
        var items = await query.OrderByDescending(q => q.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync();

        return (items, total);
    }

    public async Task<Guid> AddQuoteAsync(Quote quote)
    {
        await _context.Quotes.AddAsync(quote);
        await _context.SaveChangesAsync();
        return quote.Id;
    }

    public async Task UpdateQuoteAsync(Quote quote)
    {
        _context.Quotes.Update(quote);
        await _context.SaveChangesAsync();
    }

    public async Task<UnderwritingHealthDeclaration?> GetHealthDeclarationByIdAsync(Guid id)
    {
        return await _context.HealthDeclarations.FirstOrDefaultAsync(h => h.Id == id);
    }

    public async Task<UnderwritingHealthDeclaration?> GetHealthDeclarationByQuoteIdAsync(Guid quoteId)
    {
        return await _context.HealthDeclarations.FirstOrDefaultAsync(h => h.QuoteId == quoteId);
    }

    public async Task AddHealthDeclarationAsync(UnderwritingHealthDeclaration declaration)
    {
        await _context.HealthDeclarations.AddAsync(declaration);
        await _context.SaveChangesAsync();
    }

    public async Task UpdateHealthDeclarationAsync(UnderwritingHealthDeclaration declaration)
    {
        _context.HealthDeclarations.Update(declaration);
        await _context.SaveChangesAsync();
    }

    public async Task<UnderwritingDecision?> GetDecisionByQuoteIdAsync(Guid quoteId)
    {
        return await _context.UnderwritingDecisions.Where(d => d.QuoteId == quoteId).OrderByDescending(d => d.DecidedAt).FirstOrDefaultAsync();
    }

    public async Task AddDecisionAsync(UnderwritingDecision decision)
    {
        await _context.UnderwritingDecisions.AddAsync(decision);
        await _context.SaveChangesAsync();
    }

    public async Task UpdateDecisionAsync(UnderwritingDecision decision)
    {
        _context.UnderwritingDecisions.Update(decision);
        await _context.SaveChangesAsync();
    }

    public async Task<(List<UnderwritingDecision> Items, int TotalCount)> GetDecisionHistoryAsync(Guid quoteId)
    {
        var query = _context.UnderwritingDecisions.Where(d => d.QuoteId == quoteId);
        var total = await query.CountAsync();
        var items = await query.OrderByDescending(d => d.DecidedAt).ToListAsync();
        return (items, total);
    }

    public async Task<long> GetNextQuoteSequenceAsync()
    {
        var result = await _context.Database
            .SqlQueryRaw<long>("SELECT nextval('insurance_schema.quote_number_seq')")
            .ToListAsync();
        return result.FirstOrDefault();
    }
}
