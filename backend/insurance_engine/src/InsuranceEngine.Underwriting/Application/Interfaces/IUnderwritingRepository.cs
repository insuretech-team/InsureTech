using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Domain.Entities;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Application.Interfaces;

public interface IUnderwritingRepository
{
    Task<Quote?> GetQuoteByIdAsync(Guid id);
    Task<Quote?> GetQuoteByNumberAsync(string quoteNumber);
    Task<(List<Quote> Items, int TotalCount)> ListQuotesAsync(
        Guid? beneficiaryId, QuoteStatus? status, int page, int pageSize);
    Task<Guid> AddQuoteAsync(Quote quote);
    Task UpdateQuoteAsync(Quote quote);

    Task<UnderwritingHealthDeclaration?> GetHealthDeclarationByIdAsync(Guid id);
    Task<UnderwritingHealthDeclaration?> GetHealthDeclarationByQuoteIdAsync(Guid quoteId);
    Task AddHealthDeclarationAsync(UnderwritingHealthDeclaration declaration);
    Task UpdateHealthDeclarationAsync(UnderwritingHealthDeclaration declaration);

    Task<UnderwritingDecision?> GetDecisionByQuoteIdAsync(Guid quoteId);
    Task AddDecisionAsync(UnderwritingDecision decision);
    Task UpdateDecisionAsync(UnderwritingDecision decision);
    Task<(List<UnderwritingDecision> Items, int TotalCount)> GetDecisionHistoryAsync(Guid quoteId);
    
    Task<long> GetNextQuoteSequenceAsync();
}
